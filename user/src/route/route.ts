import express, { Router } from 'express'
import { findID, loginUser, registerUser } from '../handler/user.handler'
import { verify } from '../middleware/verify'

export const router: Router = express.Router()
router.post("/register", registerUser)
router.post("/login", loginUser)
router.get("/get/:id", verify, findID)

