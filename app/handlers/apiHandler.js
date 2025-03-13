const {publishMessage} = require("./natsHandler.js");
const {handleReadAll} = require("./postgresHandler");

const handleStartProcess = (req, res) => {
    const {filePath} = req.body;

    if (!filePath) {
        return res.status(400).json({error: 'filePath is required'});
    }
    try {
        publishMessage('mp4FilePaths', filePath);
        res.status(200).json({message: 'filePath published to NATS successfully'});
    } catch (err) {
        console.error('Error in handleStart:', err);
        res.status(500).json({error: 'Failed to publish to NATS'});
    }
};

const handleListAll = async (req, res) => {
    let allProtocols = await handleReadAll()
    res.status(200).json({message: allProtocols});
};
const handleListDetail = async (req, res) => {
};
const handleDelete = async (req, res) => {
};

module.exports = {
    handleStartProcess,
    handleListAll,
    handleListDetail,
    handleDelete
}
