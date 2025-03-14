require('dotenv').config();

const {
    NATS_URL: nc_url = "nats://localhost:4222",
    MP4_FILE_PATHS_TOPIC: mp4FilePathsTopic="mp4FilePaths",
    INITIAL_SEGMENT_FILE_PATHS: InitialSegmentFilePaths= "InitialSegmentFilePaths"

} = process.env;

module.exports = {
    nc_url,
    mp4FilePathsTopic,
    InitialSegmentFilePaths
};