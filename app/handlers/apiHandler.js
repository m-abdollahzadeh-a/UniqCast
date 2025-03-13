const {publishMessage} = require("./natsHandler.js");
const {handleReadAll, handleReadDetail, handleDeleteWithFileName} = require("./postgresHandler");

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
    const {fileName} = req.body;
    if (!fileName) {
        return res.status(400).json({error: 'fileName is required'});
    }
    try {
        let detail = await handleReadDetail(fileName)
        res.status(200).json({message: detail});
    } catch (err) {
        console.error('Error in handleListDetail:', err);
        res.status(500).json({error: 'Failed to List file process result'});
    }
};
const handleDelete = async (req, res) => {
    const {fileName} = req.body;
    if (!fileName) {
        return res.status(400).json({error: 'fileName is required'});
    }
    try {
        let deleted = await handleDeleteWithFileName(fileName)
        res.status(200).json({message: `Deleted: ${deleted}`});
    } catch (err) {
        console.error('Error in handleListDetail:', err);
        res.status(500).json({error: 'Failed to List file process result'});
    }
};

module.exports = {
    handleStartProcess,
    handleListAll,
    handleListDetail,
    handleDelete
}
