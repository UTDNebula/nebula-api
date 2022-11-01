/////////////////////////////////
//	Script for converting combined PsuedoCourses into Courses by parsing their requirements. Should be run after combining the scraped data with Combine.ts.
////////////////////////////////

import { ParsingUtils } from './Utils';
ParsingUtils.ParseAllReqs();