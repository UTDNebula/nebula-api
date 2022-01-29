import { Request, Response, NextFunction } from 'express';
import { isValidObjectId } from 'mongoose';

import { TokenModel } from "../models/token";

export default async function token(req: Request, res: Response, next: NextFunction) {
  const tokenId = req.header('Authorization');
  if (!isValidObjectId(tokenId)) {
    return res.status(400).json({
      error: 'invalid authorization token'
    });
  }
  const token = await TokenModel.findById(tokenId).exec();
  if (!token) {
    return res.status(500).json({
      error: 'authorization token not found'
    });
  }
  if (!token.active) {
    return res.status(400).json({
      error: 'authorization token is not active'
    });
  }
  next();
}