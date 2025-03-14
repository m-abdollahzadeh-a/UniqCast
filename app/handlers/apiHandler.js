const {publishMessage} = require("./natsHandler.js");
const {handleReadAll, handleReadDetail, handleDeleteWithID} = require("./postgresHandler");

const handleStartProcess = (req, res) => {
    const {filePath} = req.body;

    if (!filePath) {
        return res.status(400).json({error: 'filePath is required'});
    }
    try {
        publishMessage('mp4FilePaths', filePath);
        res.status(200).json({message: 'File started Processing'});
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
    const {id} = req.params;
    if (!id) {
        return res.status(400).json({error: 'Message Not Found'});
    }
    try {
        let detail = await handleReadDetail(id)
        if (detail === undefined) {
            res.status(404).json({message: `ID not exists`});
        }else {
            res.status(200).json({message: detail});
        }
    } catch (err) {
        console.error('Error in handleListDetail:', err);
        res.status(500).json({error: 'Failed to List file process result'});
    }
};
const handleDelete = async (req, res) => {
    const {id} = req.params;
    if (!id) {
        return res.status(400).json({error: 'ID is required'});
    }
    try {
        let deleted = await handleDeleteWithID(id)
        if (deleted === undefined) {
            res.status(404).json({message: `ID not exists`});
        }else {
            res.status(200).json({message: `Deleted`});
        }

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
