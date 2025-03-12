const express = require('express');
const { connect } = require('nats');

const app = express();
const port = 3000;

app.use(express.json());

// NATS connection
let nc;

(async () => {
    try {
        nc = await connect({ servers: 'nats://localhost:4222' }); // Replace with your NATS server URL
        console.log('Connected to NATS');

        nc.closed()
            .then(() => {
                console.log('NATS connection closed');
            })
            .catch((err) => {
                console.error('NATS connection error:', err);
            });
    } catch (err) {
        console.error('Failed to connect to NATS:', err);
        process.exit(1);
    }
})();

app.post('/start', async (req, res) => {
    const { filePath } = req.body;

    if (!filePath) {
        return res.status(400).json({ error: 'filePath is required' });
    }

    if (!nc) {
        return res.status(500).json({ error: 'NATS connection not established' });
    }

    try {
        // Publish the filePath to a NATS subject
        nc.publish('mp4FilePaths', Buffer.from(filePath));
        console.log(`filePath "${filePath}" published to NATS`);
        res.status(200).json({ message: 'filePath published to NATS successfully' });
    } catch (err) {
        console.error('Failed to publish to NATS:', err);
        res.status(500).json({ error: 'Failed to publish to NATS' });
    }
});

app.listen(port, () => {
    console.log(`Server is running on http://localhost:${port}`);
});

// Graceful shutdown
process.on('SIGINT', async () => {
    if (nc) {
        await nc.close();
        console.log('NATS connection closed');
    }
    process.exit();
});