import jsonwebtoken from "jsonwebtoken";
import "dotenv/config";

export class Payload  {
    username: string
    avatar: string
    roleID: number
}

export const generateToken = (payload: Payload) => {
    return jsonwebtoken.sign(payload, process.env.APP_SECRET, {
        expiresIn: '1d',
        algorithm: 'HS256'
    })
}
