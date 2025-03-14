const express = require('express');
const {connectToNATS, closeNATSConnection} = require('./handlers/natsHandler.js');
const {handleStartProcess, handleListAll, handleDelete, handleListDetail} = require('./handlers/apiHandler.js');
const {handleWriteToPostgres} = require('./handlers/postgresHandler.js');
const setupSwagger = require('./swagger');

require('dotenv').config();

const app = express();
setupSwagger(app);

const {PORT: port = 3000,} = process.env;

app.use(express.json());

connectToNATS()
    .then(() => {
        console.log('NATS connection established');
    })
    .catch((err) => {
        console.error('Failed to connect to NATS:', err);
        process.exit(1);
    });

/**
 * @swagger
 * /process:
 *   post:
 *     summary: Start a file processing task
 *     description: Initiates a processing task for the file specified in the input body.
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               filePath:
 *                 type: string
 *                 example: "/home/jan/Documents/video.mp4"
 *             required:
 *               - filePath
 *     responses:
 *       200:
 *         description: Processing started successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 message:
 *                   type: string
 *                   example: "Processing started"
 *       400:
 *         description: Invalid input (e.g., missing or invalid filePath)
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 message:
 *                   type: string
 *                   example: "Invalid file path"
 *       500:
 *         description: Internal server error
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 message:
 *                   type: string
 *                   example: "Internal server error"
 */
app.post('/process', handleStartProcess);

/**
 * @swagger
 * /list/all:
 *   get:
 *     summary: Retrieve a list of all items
 *     description: Returns a list of all items with their details.
 *     responses:
 *       200:
 *         description: Successful response with a list of items
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 message:
 *                   type: array
 *                   items:
 *                     type: object
 *                     properties:
 *                       id:
 *                         type: integer
 *                         example: 1
 *                       fileName:
 *                         type: string
 *                         example: "/home/jan/Documents/video.mp4"
 *                       StatusCode:
 *                         type: string
 *                         example: "Successful"
 *                       Message:
 *                         type: string
 *                         example: "File processed successfully"
 *                       ResultPath:
 *                         type: string
 *                         example: "/tmp/outputs/video.mp4"
 *                       createdAt:
 *                         type: string
 *                         format: date-time
 *                         example: "2025-03-14T15:30:30.441Z"
 *                       updatedAt:
 *                         type: string
 *                         format: date-time
 *                         example: "2025-03-14T15:30:30.441Z"
 *       500:
 *         description: Internal server error
 */
app.get('/list/all', handleListAll);

/**
 * @swagger
 * /list/detail/{id}:
 *   get:
 *     summary: Retrieve details of a specific item by ID
 *     description: Returns the details of an item based on the provided ID.
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: integer
 *         description: The ID of the item to retrieve
 *     responses:
 *       200:
 *         description: Successful response with item details
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 message:
 *                   type: object
 *                   properties:
 *                     id:
 *                       type: integer
 *                       example: 1
 *                     fileName:
 *                       type: string
 *                       example: "/home/jan/Documents/video.mp4"
 *                     StatusCode:
 *                       type: string
 *                       example: "Successful"
 *                     Message:
 *                       type: string
 *                       example: "File processed successfully"
 *                     ResultPath:
 *                       type: string
 *                       example: "/tmp/outputs/video.mp4"
 *                     createdAt:
 *                       type: string
 *                       format: date-time
 *                       example: "2025-03-14T15:30:30.441Z"
 *                     updatedAt:
 *                       type: string
 *                       format: date-time
 *                       example: "2025-03-14T15:30:30.441Z"
 *       404:
 *         description: Item not found
 *       500:
 *         description: Internal server error
 */
app.get('/list/detail/:id', handleListDetail);

/**
 * @swagger
 * /delete/{id}:
 *   delete:
 *     summary: Delete a message by ID
 *     description: Deletes a message with the specified ID. Returns a success or failure message.
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: integer
 *         description: The ID of the item to delete
 *     responses:
 *       200:
 *         description: Item successfully deleted
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 message:
 *                   type: string
 *                   example: "Deleted"
 *       404:
 *         description: Item not found
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 message:
 *                   type: string
 *                   example: "ID not exists"
 *       500:
 *         description: Internal server error
 */
app.delete('/delete/:id', handleDelete);

app.listen(port, () => {
    console.log(`Server is running on http://localhost:${port}`);
});

process.on('SIGINT', async () => {
    await closeNATSConnection();
    process.exit();
});

// Initialize Postgres table and Database
const sequelize = require('./config/database');
const Protocol = require("./models/Protocol");
const {nc_url} = require("./config/nats");
const syncDatabase = async () => {
    try {
        await sequelize.sync({force: true}); // `force: true` will drop the table if it already exists
        console.log('Database synced successfully.');
    } catch (error) {
        console.error('Error syncing database:', error);
    }
};
syncDatabase();

handleWriteToPostgres(nc_url)
