const { connect } = require('nats');

let nc;

const connectToNATS = async (nc_url) => {
    try {
        nc = await connect({ servers: nc_url });
        console.log('Connected to NATS');
        return nc;
    } catch (err) {
        console.error('Failed to connect to NATS:', err);
        throw err;
    }
};

const publishMessage = (subject, message) => {
    if (!nc) {
        throw new Error('NATS connection not established');
    }

    try {
        nc.publish(subject, Buffer.from(message));
        console.log(`Message "${message}" published to NATS subject "${subject}"`);
    } catch (err) {
        console.error('Failed to publish to NATS:', err);
        throw err;
    }
};

const closeNATSConnection = async () => {
    if (nc) {
        await nc.close();
        console.log('NATS connection closed');
    }
};

module.exports = {
    connectToNATS,
    publishMessage,
    closeNATSConnection,
};