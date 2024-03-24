import { Request, Response } from "express";
import bcrypt from 'bcrypt';
import { prisma } from "..";
import { Payload, generateToken } from "../middleware/jwt";

export class Result {
    fname: string
    lname: string
    username: string
    avatar: string
    roleID: number
}

export const registerUser = async (req: Request, res: Response) => {
    try {
        const { fname, lname, username, email, avatar, roleID, password } = req.body
        const hashPassword = bcrypt.hash(password, 14)

        const userCreate = await prisma.user.create({
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
        result.avatar = userCreate.avatar
        result.fname = userCreate.fname
        result.lname = userCreate.lname
        result.roleID = userCreate.roleID
        result.username = userCreate.username

        res.status(200).json({
            msg: 'success creata account',
            code: 200,
            data: result,
        })

    } catch (error) {
        res.status(500).json({
            msg: 'something error',
            err: error,
        })
    }
}

export const loginUser = async (req: Request, res: Response) => {
    try {
        const { email, password } = req.body
        const find = await prisma.user.findUnique({
            where: {
                email: email
            }
        })

        if (!find) {
            res.status(404).json({ msg: 'sorry your email not found', code: 404 })
        }

        const verify = bcrypt.compare(password, find.password)
        if (!verify) {
            res.status(400).json({
                msg: 'sorry your password not matched',
                code: 400
            })
        }

        const payload = new Payload()
        payload.avatar = find.avatar
        payload.roleID = find.roleID
        payload.username = find.username

        const token = await generateToken(payload)

        res.cookie('authorization', token, {
            httpOnly: true,
            sameSite: "lax"
        })

        res.status(200).json({ msg: 'success login' })
    } catch (error) {
        res.status(500).json({
            msg: 'something error',
            err: error,
        })
    }
}

export const findID = async (req: Request, res: Response) => {
    try {
        const { id } = req.params
        const result = await prisma.user.findUnique({
            where: {
                id: Number(id)
            }
        })

        const response: Result = new Result()
        response.avatar = result.avatar
        response.fname = result.fname
        response.lname = result.lname
        response.roleID = result.roleID
        response.username = result.username

        if (!result) {
            res.status(404).json({
                msg: 'user not found'
            })
        }

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