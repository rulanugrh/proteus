import { PrismaClient } from '@prisma/client';
import express, { Request, Response } from 'express'
import 'dotenv/config'
import { router } from "./route/route";
import { register, totalCPU, totalMemory } from "./helper/prometheus";
import { Logger } from "tslog";
import { appendFileSync } from "fs";

// setup logger for trigger to loki
// and save logger tu local
export const logger = new Logger({ name: 'user-services '});
logger.attachTransport((logger) => {
    appendFileSync("../../data/log/user.log", JSON.stringify(logger) + '\n')
})

export const prisma = new PrismaClient();

const app = express()
const port = process.env.APP_PORT

// setup for cpu usage server
const cpuUsage = process.cpuUsage()
const currentUsage = (cpuUsage.user + cpuUsage.system) * 1000

// get total memory usage on server
const mem = process.memoryUsage()

async function main() {
    totalCPU.labels('v1').set(currentUsage)
    totalMemory.labels('v1').set(mem.heapUsed)

    app.use(express.json())

    app.use("/api/user", router)

    app.get('/metrics', async (req: Request, res: Response) => {
        res.setHeader('Content-Type', register.contentType)
        res.send(await register.metrics())
    })

    app.all("*", (req: Request, res: Response) => {
        res.status(404).json({ error: `Route ${req.originalUrl} not found` });
    });

    app.listen(port, () => {
        console.log(`Server running at: ${port}`)
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