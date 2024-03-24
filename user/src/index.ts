import { PrismaClient } from "@prisma/client";
import express, { Request, Response } from 'express'
import 'dotenv/config'
import { router } from "./route/route";

export const prisma = new PrismaClient();


const app = express()
const port = process.env.APP_PORT

async function main() {
    app.use(express.json())

    app.use("/api/user", router)

    app.all("*", (req: Request, res: Response) => {
        res.status(404).json({ error: `Route ${req.originalUrl} not found` });
    });

    app.listen(process.env.APP_PORT, () => {
        console.log(`Server running at: ${process.env.APP_PORT}`)
    })
}

main()
    .then(async () => {
        await prisma.$connect()
    })
    .catch(async (e) => {
        console.error(e)
        await prisma.$disconnect()
        process.exit(1)
    })