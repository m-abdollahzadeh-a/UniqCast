const { DataTypes } = require('sequelize');
const sequelize = require('../config/database');

const Protocol = sequelize.define('Protocol', {
    fileName: {
        type: DataTypes.STRING,
    },
    StatusCode: {
        type: DataTypes.STRING,
    },
    Message: {
        type: DataTypes.STRING,
    },
    ResultPath: {
        type: DataTypes.STRING,
    },
});

module.exports = Protocol;