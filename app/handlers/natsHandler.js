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

const subscribeToSubject = (subject, callback) => {
    if (!nc) {
        throw new Error('NATS connection not established');
    }

    try {
        const subscription = nc.subscribe(subject, {
            callback: (err, msg) => {
                if (err) {
                    console.error('Error processing message:', err);
                    return;
                }
                callback(msg);
            }
        });

        console.log(`Subscribed to NATS subject "${subject}"`);
        return subscription;
    } catch (err) {
        console.error('Failed to subscribe to NATS:', err);
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
    closeNATSConnection,
    subscribeToSubject,
    publishMessage,
    connectToNATS
}
