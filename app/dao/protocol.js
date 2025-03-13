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

const findDetailProtocol = async (file_name) => {
    try {
        const protocol = await Protocol.findOne({
            where: {fileName: file_name},
        });

        if (protocol) {
            console.log('protocol found:', protocol.toJSON());
            return protocol.toJSON()
        } else {
            console.log('protocol not found');
        }
    } catch (error) {
        console.error('Error finding protocol:', error);
    }
};

module.exports = {
    createProtocol,
    findAllProtocols,
    findDetailProtocol
}
