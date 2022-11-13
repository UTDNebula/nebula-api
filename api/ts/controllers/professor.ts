import { Request, Response } from 'express';
import { ProfessorModel } from '../models/professor';

export const professorSearch = async (req: Request, res: Response) => {
  if (Object.keys(req.query).length < 1) {
    return res.status(400).json({
      error: 'request did not contain any query parameters',
    });
  }
  ProfessorModel.find(req.query, {}, { strict: false }, (error, result) => {
    if (error) {
      return res.status(500).json({
        error: 'an internal server error occured',
      });
    } else {
      return res.status(200).json(result);
    }
  });
};

export const professorById = async (req: Request, res: Response) => {
  if (req.params.id === null) {
    res.status(400).json({
      error: 'request did not contain a professor id',
    });
  }
  ProfessorModel.findOne({ _id: req.params.id }, {}, {}, (error, result) => {
    if (error) {
      res.status(500).json({
        error: 'an internal server error occured',
      });
    } else {
      res.status(200).json(result);
    }
  });
};
