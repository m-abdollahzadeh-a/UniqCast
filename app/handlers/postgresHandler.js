const { subscribeToSubject } = require ('./natsHandler.js');
const {closeNATSConnection, connectToNATS} = require("./natsHandler");

const handleMessage = (msg) => {
    console.log(`Received message: ${msg.data.toString()}`);
    // to postgres

};

const writeResultToPostgres = async () => {
    try {
        await connectToNATS('nats://localhost:4222');
        subscribeToSubject('InitialSegmentFilePaths', handleMessage);
        process.on('SIGINT', async () => {
            await closeNATSConnection();
            process.exit();
        });
    } catch (err) {
        console.error('Error:', err);
    }
};

const handleWriteToPostgres = async () => {
    console.log('Starting NATS listener...');
    await writeResultToPostgres();
};


module.exports = {
    handleWriteToPostgres
}
