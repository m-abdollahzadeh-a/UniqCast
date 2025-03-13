const express = require('express');
const {connectToNATS, closeNATSConnection} = require('./handlers/natsHandler');
const {handleStartProcess, handleListAll, handleDelete, handleListDetail} = require('./handlers/apiHandler');

const app = express();
const port = 3000;
const nc_url = "nats://localhost:4222"

app.use(express.json());

connectToNATS(nc_url)
    .then(() => {
        console.log('NATS connection established');
    })
    .catch((err) => {
        console.error('Failed to connect to NATS:', err);
        process.exit(1);
    });

app.post('/process', handleStartProcess);
app.get('/list/all', handleListAll);
app.post('/list/detail', handleListDetail);
app.delete('/delete', handleDelete);

app.listen(port, () => {
    console.log(`Server is running on http://localhost:${port}`);
});

process.on('SIGINT', async () => {
    await closeNATSConnection();
    process.exit();
});
