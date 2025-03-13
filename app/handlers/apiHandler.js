const {publishMessage} = require('./natsHandler');

const handleStartProcess = async (req, res) => {
    const {filePath} = req.body;

    if (!filePath) {
        return res.status(400).json({error: 'filePath is required'});
    }

    try {
        await publishMessage('mp4FilePaths', filePath);
        res.status(200).json({message: 'filePath published to NATS successfully'});
    } catch (err) {
        console.error('Error in handleStart:', err);
        res.status(500).json({error: 'Failed to publish to NATS'});
    }
};

module.exports = {
    handleStartProcess,
};
