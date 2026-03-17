import React, { useState } from 'react';
import { callApi, getUserIdFromJwt } from '../utils/api';

const Agent = ({ activeModel, setApiResponse }) => {
    const [ticketDetails, setTicketDetails] = useState(null);
    const [ticketComments, setTicketComments] = useState([]);
    const [commentSearchTerm, setCommentSearchTerm] = useState('');

    const handleNumericInput = (e) => {
        e.target.value = e.target.value.replace(/[^0-9]/g, '');
    };

    const handleNumericKeyDown = (e) => {
        let currentValue = parseInt(e.target.value) || 0;
        if (e.key === "ArrowUp") {
            e.preventDefault();
            e.target.value = currentValue + 1;
        } else if (e.key === "ArrowDown") {
            e.preventDefault();
            e.target.value = Math.max(0, currentValue - 1);
        }
    };

    const handleLoadAgentTicket = async (e) => {
        e.preventDefault();
        const ticketId = e.target.agentTicketId.value;
        if (!ticketId) return;

        const data = await callApi("GET", `/agent/tickets/${ticketId}`, null, setApiResponse);
        if (data) {
            setTicketDetails(data);
            setTicketComments(data.Comments || []);
        } else {
            setTicketDetails(null);
            setTicketComments([]);
        }
    };

    const handleUpdateAgentTicket = async (e) => {
        e.preventDefault();
        const ticketId = ticketDetails.ID;
        const status = e.target.agentTicketStatus.value;
        const priority = e.target.agentTicketPriority.value;

        const ticketData = {
            status: status,
            priority: priority,
        };

        await callApi("PUT", `/agent/tickets/${ticketId}`, ticketData, setApiResponse);
    };

    const handleAddAgentComment = async (e) => {
        e.preventDefault();
        const ticketId = parseInt(e.target.agentCommentTicketId.value);
        const body = e.target.agentCommentBody.value;
        const isInternal = e.target.agentCommentIsInternal.checked;

        const commentData = {
            ticket_id: ticketId,
            body: body,
            is_internal: isInternal,
        };

        await callApi("POST", `/agent/tickets/${ticketId}/comments`, commentData, setApiResponse);
        if (ticketDetails) {
            const reloadedData = await callApi("GET", `/agent/tickets/${ticketDetails.ID}`, null, setApiResponse);
            if (reloadedData) {
                setTicketDetails(reloadedData);
                setTicketComments(reloadedData.Comments || []);
            }
        }
    };

    const handleAssignAgentTicket = async (e) => {
        e.preventDefault();
        const ticketId = e.target.assignTicketId.value;
        if (!ticketId) return;

        const userId = getUserIdFromJwt();
        if (userId === null) {
            alert("Could not get user ID from JWT. Please log in again.");
            return;
        }
        const assigneeIdInt = parseInt(userId);
        await callApi("PUT", `/agent/tickets/${ticketId}`, { AssigneeID: assigneeIdInt }, setApiResponse);
    };

    const filteredComments = ticketComments.filter(comment =>
        comment.Body?.toLowerCase().includes(commentSearchTerm.toLowerCase())
    );

    if (activeModel !== 'agent') {
        return null;
    }

    return (
        <div id="agent" className="model-section">
            <h2>Agent Features</h2>
            <div>
                <h3>Agent Ticket Actions</h3>
                <div className="button-group">
                    <button onClick={() => callApi('GET', `/agent/tickets`, null, setApiResponse)}>List My Tickets</button>
                    <button onClick={() => callApi('GET', `/agent/tickets/open`, null, setApiResponse)}>List All Open Tickets</button>
                </div>
            </div>

            <div>
                <h3>View/Update Assigned Ticket (GET/PUT /agent/tickets/{'{id}'})</h3>
                <form id="loadAgentTicketForm" onSubmit={handleLoadAgentTicket}>
                    <label htmlFor="agentTicketId">Ticket ID:</label>
                    <input type="text" id="agentTicketId" name="agentTicketId" required onInput={handleNumericInput} onKeyDown={handleNumericKeyDown} />
                    <button type="submit">Load Ticket</button>
                </form>
                {ticketDetails && (
                    <form id="updateAgentTicketForm" style={{ paddingTop: '5px' }} onSubmit={handleUpdateAgentTicket}>
                        <input type="hidden" id="hiddenAgentTicketId" name="hiddenAgentTicketId" value={ticketDetails.ID} />
                        <div className="form-group-item">
                            <label htmlFor="agentTicketTitle">Title:</label>
                            <input type="text" id="agentTicketTitle" name="agentTicketTitle" readOnly value={ticketDetails?.Title || ''} />
                        </div>
                        <div className="form-group-item full-width-item">
                            <label htmlFor="agentTicketDescription">Description:</label>
                            <textarea id="agentTicketDescription" name="agentTicketDescription" readOnly value={ticketDetails?.Description || ''}></textarea>
                        </div>
                        <div className="form-group-grid">
                            <div className="form-group-item">
                                <label htmlFor="agentTicketStatus">Status:</label>
                                <select id="agentTicketStatus" name="agentTicketStatus" value={ticketDetails?.Status || 'Open'} onChange={(e) => setTicketDetails({ ...ticketDetails, Status: e.target.value })}>
                                    <option value="Open">Open</option>
                                    <option value="In Progress">In Progress</option>
                                    <option value="Resolved">Resolved</option>
                                    <option value="Closed">Closed</option>
                                </select>
                            </div>
                            <div className="form-group-item">
                                <label htmlFor="agentTicketPriority">Priority:</label>
                                <select id="agentTicketPriority" name="agentTicketPriority" value={ticketDetails?.Priority || 'Low'} onChange={(e) => setTicketDetails({ ...ticketDetails, Priority: e.target.value })}>
                                    <option value="Low">Low</option>
                                    <option value="Medium">Medium</option>
                                    <option value="High">High</option>
                                    <option value="Urgent">Urgent</option>
                                </select>
                            </div>
                        </div>
                        <button type="submit">Update Ticket</button>
                    </form>
                )}
            </div>

            <div>
                <h3>Add Comment to Ticket (POST /agent/tickets/{'{id}'}/comments)</h3>
                <form id="addAgentCommentForm" onSubmit={handleAddAgentComment}>
                    <div className="form-group-item">
                        <label htmlFor="agentCommentTicketId">Ticket ID:</label>
                        <input type="text" id="agentCommentTicketId" name="agentCommentTicketId" required onInput={handleNumericInput} onKeyDown={handleNumericKeyDown} value={ticketDetails?.ID || ''} readOnly={!!ticketDetails} />
                    </div>
                    <div className="form-group-item full-width-item">
                        <label htmlFor="agentCommentBody">Comment:</label>
                        <textarea id="agentCommentBody" name="agentCommentBody" required></textarea>
                    </div>
                    <div className="form-group-item">
                        <label htmlFor="agentCommentIsInternal">Internal Only:</label>
                        <input type="checkbox" id="agentCommentIsInternal" name="agentCommentIsInternal" style={{ maxWidth: '25%', minHeight: '42px' }} />
                    </div>
                    <button type="submit">Add Comment</button>
                </form>
            </div>

            {ticketComments.length > 0 && (
                <div id="agentTicketCommentsSection">
                    <h3>Comments for this Ticket</h3>
                    <input type="text" id="agentCommentSearchInput" placeholder="Search comments..." style={{ width: '90%', marginBottom: '10px', marginLeft: '3%' }} onChange={(e) => setCommentSearchTerm(e.target.value)} />
                    <div style={{ maxHeight: '600px', overflowY: 'auto', border: '1px solid #0f0' }}>
                        <table id="agentTicketCommentsTable" style={{ width: '100%', marginTop: '10px' }}>
                            <thead>
                                <tr>
                                    <th style={{ border: '1px solid #0f0', padding: '8px', textAlign: 'center', position: 'sticky', top: 0, backgroundColor: '#000', zIndex: 1 }}>ID</th>
                                    <th style={{ border: '1px solid #0f0', padding: '8px', textAlign: 'center', position: 'sticky', top: 0, backgroundColor: '#000', zIndex: 1 }}>Author ID</th>
                                    <th style={{ border: '1px solid #0f0', padding: '8px', textAlign: 'left', position: 'sticky', top: 0, backgroundColor: '#000', zIndex: 1 }}>Body</th>
                                    <th style={{ border: '1px solid #0f0', padding: '8px', textAlign: 'center', position: 'sticky', top: 0, backgroundColor: '#000', zIndex: 1 }}>Internal</th>
                                    <th style={{ border: '1px solid #0f0', padding: '8px', textAlign: 'left', position: 'sticky', top: 0, backgroundColor: '#000', zIndex: 1 }}>Created At</th>
                                </tr>
                            </thead>
                            <tbody>
                                {filteredComments.map(comment => (
                                    <tr key={comment.ID}>
                                        <td style={{ border: '1px solid #0f0', textAlign: 'center' }}>{comment.ID}</td>
                                        <td style={{ border: '1px solid #0f0', textAlign: 'center' }}>{comment.AuthorID}</td>
                                        <td style={{ border: '1px solid #0f0' }}>{comment.Body}</td>
                                        <td style={{ border: '1px solid #0f0', textAlign: 'center' }}>{comment.IsInternal ? "Yes" : "No"}</td>
                                        <td style={{ border: '1px solid #0f0', textAlign: 'left' }}>{new Date(comment.CreatedAt).toLocaleString()}</td>
                                    </tr>
                                ))}
                                {filteredComments.length === 0 && (
                                    <tr>
                                        <td colSpan="5" style={{ border: '1px solid #0f0' }}>No comments found for this ticket.</td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>
            )}

            <div>
                <h3>Assign Ticket to Self (PUT /agent/tickets/{'{id}'})</h3>
                <form id="assignAgentTicketForm" onSubmit={handleAssignAgentTicket}>
                    <div className="form-group-item">
                        <label htmlFor="assignTicketId">Ticket ID:</label>
                        <input type="text" id="assignTicketId" name="assignTicketId" required onInput={handleNumericInput} onKeyDown={handleNumericKeyDown} />
                    </div>
                    <button type="submit">Assign to Self</button>
                </form>
            </div>
        </div>
    );
};

export default Agent;