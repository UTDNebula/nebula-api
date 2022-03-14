import schemas from '../api/schemas';
import mongoose from 'mongoose';
import { readFileSync } from 'fs';

const Professors = mongoose.model('Professors', new mongoose.Schema<schemas.Professor>());

async function main() : Promise<void> {
    await mongoose.connect(`mongodb+srv://${process.env.MONGO_USERNAME}:${process.env.MONGO_PASSWORD}@development-0.gftz1.mongodb.net/professorDB?retryWrites=true&w=majority`);

    let rawdata = readFileSync('./data/Professors.json');
    let professors = JSON.parse(rawdata.toString());

    Professors.insertMany(professors).then(() => {
        console.log('done (with or without errors)');
    });
}

main();
