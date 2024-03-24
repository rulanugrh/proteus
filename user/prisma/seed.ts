import { PrismaClient } from "@prisma/client";

const prisma = new PrismaClient();

const roles = [
    {
        name: "Administrator",
        desc: "This is role administrator"
    },
    {
        name: "Owner",
        desc: "This is role owner"
    }
]

async function main() {
    console.log("Start seeder to Database")
    for (const r of roles) {
        const role = await prisma.role.create({
            data: r
        })
        console.log(`Create role by id: ${role.id}`)
    }

    console.log(`Seed finish..`)
}

main()
    .then(async() => {
        await prisma.$disconnect()
    })
    .catch(async(e) => {
        console.error(e)
        await prisma.$disconnect()
        process.exit(1)
    })