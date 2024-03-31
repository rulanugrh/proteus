import express, { Router } from 'express'
import { findID, loginUser, registerUser } from '../handler/user.handler'
import { verify } from '../middleware/verify'
import { validate } from '../middleware/validate'
import { schemaLogin, schemaRegiter } from '../helper/schema'

export const router: Router = express.Router()
router.post("/register", validate(schemaRegiter), registerUser)
router.post("/login", validate(schemaLogin), loginUser)
router.get("/get/:id", verify, findID)