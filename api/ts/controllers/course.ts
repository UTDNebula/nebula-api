import { Request, Response } from 'express';
import { CourseModel } from '../models/course';

export const courseSearch = async (req: Request, res: Response) => {
  if (Object.keys(req.query).length < 1) {
    return res.status(400).json({
      error: 'request did not contain any query parameters',
    });
  }
  CourseModel.find(req.query, {}, { strict: false }, (error, result) => {
    if (error) {
      return res.status(500).json({
        error: 'an internal server error occured',
      });
    } else {
      return res.status(200).json(result);
    }
  });
};

export const courseById = async (req: Request, res: Response) => {
  if (req.params.id === null) {
    res.status(400).json({
      error: 'request did not contain a course id',
    });
  }

  CourseModel.findOne({ _id: req.params.id }, {}, {}, (error, result) => {
    if (error) {
      res.status(500).json({
        error: 'an internal server error occured',
      });
    } else {
      res.status(200).json(result);
    }
  });
};
