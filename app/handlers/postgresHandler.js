const {subscribeToSubject} = require('./natsHandler.js');
const {closeNATSConnection, connectToNATS} = require("./natsHandler");
const {createProtocol, findAllProtocols, findDetailProtocol, deleteProtocol} = require("../dao/protocol");


const handleMessage = async (msg) => {
    console.log(`Received message: ${msg.data.toString()}`);

    const jsonMessage = JSON.parse(msg.data);
    await createProtocol(jsonMessage.file_name, jsonMessage.status_code, jsonMessage.Message, jsonMessage.result_path);
};

const writeResultToPostgres = async (nats_url) => {
    try {
        await connectToNATS(nats_url);
        subscribeToSubject('InitialSegmentFilePaths', handleMessage);
        process.on('SIGINT', async () => {
            await closeNATSConnection();
            process.exit();
        });
    } catch (err) {
        console.error('Error:', err);
    }
};

const handleWriteToPostgres = async (nats_url) => {
    console.log('Starting NATS listener...');
    await writeResultToPostgres(nats_url);
};

const handleReadAll = async () => {
    console.log('Finding All Protocols');
    return await findAllProtocols()
}

const handleReadDetail = async (id) => {
    console.log('Finding file Message');
    return await findDetailProtocol(id)
}

const handleDeleteWithID = async (id) => {
    console.log('Deleting message');
    return await deleteProtocol(id)
}

module.exports = {
    handleWriteToPostgres,
    handleReadAll,
    handleReadDetail,
    handleDeleteWithID
}
