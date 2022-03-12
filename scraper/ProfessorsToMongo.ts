import mongoose from 'mongoose';
import { readFileSync } from 'fs';

const locationSchema = new mongoose.Schema({
    building: String,
    room: String,
    map_uri: String,
});

const meetingSchema = new mongoose.Schema({
    start_date: String,
    end_date: String,
    meeting_days: [String],
    start_time: String,
    end_time: String,
    modality: String,
    location: locationSchema,
  });

const professorSchema = new mongoose.Schema({
    first_name: String,
    last_name: String,
    titles: [String],
    email: String,
    phone_number: String,
    office: locationSchema,
    profile_uri: String,
    image_uri: String,
    office_hours: [meetingSchema],
    sections: [mongoose.Types.ObjectId],
});

const Professors = mongoose.model('Professors', professorSchema);

async function main() : Promise<void> {
    await mongoose.connect(`mongodb+srv://${process.env.MONGO_USERNAME}:${process.env.MONGO_PASSWORD}@development-0.gftz1.mongodb.net/professorDB?retryWrites=true&w=majority`);

    let rawdata = readFileSync('./scraper/data/Professors.json');
    let professors = JSON.parse(rawdata.toString());

    Professors.insertMany(professors).then(() => {
        console.log('done (with or without errors)');
    });
}

main();
