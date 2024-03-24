import jwt from 'jsonwebtoken'
import { NextFunction, Request, Response } from 'express'

export const verify = async (req: Request, res: Response, next: NextFunction) => {
    const headers = req.headers["authorization"]
    const token = headers && headers.split(' ')[1]

    if (token === null ){
        return res.status(401).json({
            msg: 'sorry you not have token',
            code: 401
        })
    }

    jwt.verify(token, process.env.APP_SECRET, (err, email) => {
        console.error(err)
        if (err) {
            return res.status(403).json({
                msg: 'sorry yout not allow login this path',
                code: 403
            })
        }

        req.body.email = email
        next()
    })
}