import { Request, Response } from 'express';
import { DegreeModel } from '../models/degree';

export const degreeSearch = async (req: Request, res: Response) => {
  if (Object.keys(req.query).length < 1) {
    return res.status(400).json({
      error: 'request did not contain any query parameters',
    });
  }
  DegreeModel.find(req.query, {}, { strict: false }, (error, result) => {
    if (error) {
      return res.status(500).json({
        error: 'an internal server error occured',
      });
    } else {
      return res.status(200).json(result);
    }
  });
};

export const degreeById = async (req: Request, res: Response) => {
  if (req.params.id === null) {
    res.status(400).json({
      error: 'request did not contain a course id',
    });
  }

  DegreeModel.findOne({ _id: req.params.id }, {}, {}, (error, result) => {
    if (error) {
      res.status(500).json({
        error: 'an internal server error occured',
      });
    } else {
      res.status(200).json(result);
    }
  });
};

export const degreeInsert = async (req: Request, res: Response) => {
  const newDegree = new DegreeModel(req.body);
  newDegree.validate((err) => {
    if (err) {
      res.status(400).json(err);
    } else {
      newDegree.save((err) => {
        if (err) {
          res.status(500).json(err);
        } else {
          res.status(200).json(newDegree);
        }
      });
    }
  });
};
