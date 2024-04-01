import { Request, Response } from "express";
import bcrypt from 'bcrypt';
import { logger, prisma } from "..";
import { Payload, generateToken } from "../middleware/jwt";
import { counter, histogram } from "../helper/prometheus";

export class Result {
    fname: string
    lname: string
    username: string
    avatar: string
    roleID: number
}

export const registerUser = async (req: Request, res: Response) => {
    const start = Date.now()
    try {
        const { fname, lname, username, email, avatar, roleID, password } = req.body
        const hashPassword = bcrypt.hash(password, 14)
        
        const responseTime = Date.now() - start
        const find = await prisma.user.findUnique({
            where: {
                email: email
            }
        })

        if (find) {
            logger.error("[register] - sorry your email have been created")
            histogram.labels(req.method, '400', 'register').observe(responseTime)
            res.status(400).json({ msg: 'sorry your email have been created', code: 400 })
        }

        await prisma.user.create({
            data: {
                fname,
                lname,
                username,
                email,
                avatar,
                roleID,
                hashPassword
            }
        })

        const result: Result = new Result()
        result.avatar = avatar
        result.fname = fname
        result.lname = lname
        result.roleID = roleID
        result.username = username

        logger.info("[register] - success regiter account")
        histogram.labels(req.method, '200', 'register').observe(responseTime)
        counter.labels('register').inc()
        res.status(200).json({
            msg: 'success create account',
            code: 200,
            data: result,
        })

    } catch (error) {
        const responseTime = Date.now() - start
        logger.fatal("[register] - endpoint user register error")
        histogram.labels(req.method, '500', 'register').observe(responseTime)
        res.status(500).json({
            msg: 'something error',
            err: error,
        })
    }
}

export const loginUser = async (req: Request, res: Response) => {
    const start = Date.now()
    try {
        const responseTime = Date.now() - start
        const { email, password } = req.body
        const find = await prisma.user.findUnique({
            where: {
                email: email
            }
        })

        if (!find) {
            logger.error("[login] - sorry your email not found")
            histogram.labels(req.method, '404', 'login').observe(responseTime)
            res.status(404).json({ msg: 'sorry your email not found', code: 404 })
        }

        const verify = bcrypt.compare(password, find[0].password)
        if (!verify) {
            logger.error("[login] - sorry your password not matched")
            histogram.labels(req.method, '400', 'login').observe(responseTime)
            res.status(400).json({
                msg: 'sorry your password not matched',
                code: 400
            })
        }

        const payload = new Payload()
        payload.avatar = find[0].avatar
        payload.roleID = find[0].roleID
        payload.username = find[0].username

        const token = await generateToken(payload)

        res.cookie('authorization', token, {
            httpOnly: true,
            sameSite: "lax",
            maxAge: 24 * 60 * 60 * 1000
        })

        logger.info("[login] - success login account")
        histogram.labels(req.method, '200', 'login').observe(responseTime)
        counter.labels('login').inc()
        res.status(200).json({ msg: 'success login' })
    } catch (error) {
        const responseTime = Date.now() - start
        histogram.labels(req.method, '500', 'login').observe(responseTime)
        logger.fatal("[login] - sorry something error in endpoint user")
        res.status(500).json({
            msg: 'something error',
            err: error,
        })
    }
}

export const findID = async (req: Request, res: Response) => {
    const start = Date.now()
    try {
        const responseTime = Date.now() - start
        const { id } = req.params
        const result = await prisma.user.findUnique({
            where: {
                id: Number(id)
            }
        })
        
        if (!result) {
            logger.error("[findID] - user not found")
            histogram.labels(req.method, '404', 'findID').observe(responseTime)
            res.status(404).json({
                msg: 'user not found'
            })
        }

        const response: Result = new Result()
        response.avatar = result[0].avatar
        response.fname = result[0].fname
        response.lname = result[0].lname
        response.roleID = result[0].roleID
        response.username = result[0].username


        logger.info("[findID] -  user found")
        histogram.labels(req.method, '200', 'findID').observe(responseTime)
        res.status(200).json({
            msg: 'user found',
            code: 200,
            data: response
        })
        
    } catch (error) {
        res.status(500).json({
            msg: 'something error',
            err: error,
        })
    }
}
