const Protocol = require("../models/Protocol");

const createProtocol = async (filename, status_code, message, result_path) => {
    try {
        const protocol = await Protocol.create({
            fileName: filename,
            StatusCode: status_code,
            Message: message,
            ResultPath: result_path,
        });
        console.log('protocol created:', protocol.toJSON());
    } catch (error) {
        console.error('Error creating protocol:', error);
    }
};

const findAllProtocols = async () => {
    try {
        const protocols = await Protocol.findAll();
        console.log('All protocols:', JSON.stringify(protocols, null, 2));
        return protocols
    } catch (error) {
        console.error('Error finding protocols:', error);
    }
};

const findDetailProtocol = async (id) => {
    try {
        const protocol = await Protocol.findOne({
            where: {id: id},
        });

        if (protocol) {
            console.log('Message found:', protocol.toJSON());
            return protocol.toJSON()
        } else {
            console.log('Message not found');
        }
    } catch (error) {
        console.error('Error finding message:', error);
    }
};

const deleteProtocol = async (id) => {
    try {
        const deletedProtocol = await Protocol.destroy({
            where: { id: id },
        });

        if (deletedProtocol) {
            console.log(`Deleted ${deletedProtocol} Message(s) with the ID ${deletedProtocol}`);
            return deletedProtocol
        } else {
            console.log('Message not found');
        }
    } catch (error) {
        console.error('Error deleting message:', error);
    }
};

module.exports = {
    createProtocol,
    findAllProtocols,
    findDetailProtocol,
    deleteProtocol
}
