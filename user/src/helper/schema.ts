import { z } from 'zod'

export const schemaRegiter = z.object({
    fname: z.string(),
    lname: z.string(),
    username: z.string(),
    email: z.string().email(),
    password: z.string().min(8),
    avatar: z.string(),
    roleID: z.number()
})

export const schemaLogin = z.object({
    email: z.string().email(),
    password: z.string().min(8)
})