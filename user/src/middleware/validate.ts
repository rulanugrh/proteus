import { Request, Response, NextFunction } from 'express'
import { z, ZodError } from 'zod'

export function validate (schema: z.ZodObject<any, any>) {
    return (req: Request, res: Response, next: NextFunction) => {
        try {
            schema.parse(req.body)
            next()
        } catch (error) {
            if ( error instanceof ZodError ) {
                const errMessage = error.errors.map((issue: any) => ({
                    message: `${issue.path.join('.')} is ${issue.message}`,
                }))
                res.status(400).json({ error: 'Invalid data', details: errMessage})
            } else {
                res.status(500).json({ error: 'Internal Server Error' })
            }
        }
    }
}