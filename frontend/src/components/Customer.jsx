import React, { useState } from 'react';
import { callApi } from '../utils/api';

const Customer = ({ activeModel, setApiResponse }) => {
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

    const handleCreateCustomerTicket = async (e) => {
        e.preventDefault();
        const title = e.target.customerTicketTitle.value;
        const description = e.target.customerTicketDescription.value;
        const priority = e.target.customerTicketPriority.value;

        const ticketData = {
            title: title,
            description: description,
            priority: priority,
        };

        await callApi("POST", `/customer/tickets`, ticketData, setApiResponse);
    };

    const handleListCustomerTickets = (e) => {
        e.preventDefault();
        callApi("GET", `/customer/tickets`, null, setApiResponse);
    };

    const handleGetCustomerTicket = async (e) => {
        e.preventDefault();
        const ticketId = e.target.customerGetTicketId.value;
        if (!ticketId) return;

        const data = await callApi("GET", `/customer/tickets/${ticketId}`, null, setApiResponse);
        if (data) {
            setTicketDetails(data);
            setTicketComments(data.Comments || []);
        } else {
            setTicketDetails(null);
            setTicketComments([]);
        }
    };

    const handleAddCustomerComment = async (e) => {
        e.preventDefault();
        const ticketId = parseInt(e.target.customerCommentTicketId.value);
        const body = e.target.customerCommentBody.value;

        const commentData = {
            ticket_id: ticketId,
            body: body,
        };

        await callApi("POST", `/customer/tickets/${ticketId}/comments`, commentData, setApiResponse);
        if (ticketDetails) {
            const reloadedData = await callApi("GET", `/customer/tickets/${ticketDetails.ID}`, null, setApiResponse);
            if (reloadedData) {
                setTicketDetails(reloadedData);
                setTicketComments(reloadedData.Comments || []);
            }
        }
    };

    const handleCloseCustomerTicket = async (e) => {
        e.preventDefault();
        const ticketId = e.target.customerCloseTicketId.value;
        await callApi("PUT", `/customer/tickets/${ticketId}`, { status: "Closed" }, setApiResponse);
    };

    const filteredComments = ticketComments.filter(comment =>
        !comment.IsInternal && comment.Body?.toLowerCase().includes(commentSearchTerm.toLowerCase())
    );

    if (activeModel !== 'customer') {
        return null;
    }

    return (
        <div id="customer" className="model-section">
            <h2>Customer Features</h2>
            <div>
                <h3>Create New Ticket (POST /customer/tickets)</h3>
                <form onSubmit={handleCreateCustomerTicket}>
                    <div className="form-group-item">
                        <label htmlFor="customerTicketTitle">Title:</label>
                        <input type="text" id="customerTicketTitle" name="customerTicketTitle" required />
                    </div>
                    <div className="form-group-item full-width-item">
                        <label htmlFor="customerTicketDescription">Description:</label>
                        <textarea id="customerTicketDescription" name="customerTicketDescription"></textarea>
                    </div>
                    <div className="form-group-item">
                        <label htmlFor="customerTicketPriority">Priority:</label>
                        <select id="customerTicketPriority" name="customerTicketPriority">
                            <option value="Low">Low</option>
                            <option value="Medium">Medium</option>
                            <option value="High">High</option>
                            <option value="Urgent">Urgent</option>
                        </select>
                    </div>
                    <button type="submit">Create Ticket</button>
                </form>
            </div>

            <div>
                <h3>List My Tickets (GET /customer/tickets)</h3>
                <form onSubmit={handleListCustomerTickets}>
                    <button type="submit">List My Tickets</button>
                </form>
            </div>

            <div>
                <h3>View My Ticket (GET /customer/tickets/{'{id}'})</h3>
                <form id="getCustomerTicketForm" onSubmit={handleGetCustomerTicket}>
                    <label htmlFor="customerGetTicketId">Ticket ID:</label>
                    <input type="text" id="customerGetTicketId" name="customerGetTicketId" required onInput={handleNumericInput} onKeyDown={handleNumericKeyDown} />
                    <button type="submit">Get My Ticket</button>
                </form>
                {ticketDetails && (
                    <div style={{ paddingTop: '5px', display: 'grid' }}>
                        <div className="form-group-item">
                            <label>ID:</label>
                            <input type="text" readOnly value={ticketDetails?.ID || ''} />
                        </div>
                        <div className="form-group-item">
                            <label>Title:</label>
                            <input type="text" readOnly value={ticketDetails?.Title || ''} />
                        </div>
                        <div className="form-group-item full-width-item">
                            <label>Description:</label>
                            <textarea readOnly value={ticketDetails?.Description || ''}></textarea>
                        </div>

                        <div className="form-group-grid">
                            <div className="form-group-item">
                                <label>Status:</label>
                                <input type="text" readOnly value={ticketDetails?.Status || ''} />
                            </div>
                            <div className="form-group-item">
                                <label>Priority:</label>
                            <input type="text" readOnly value={ticketDetails?.Priority || ''} />
                            </div>
                        </div>

                        <div className="form-group-grid">
                            <div className="form-group-item">
                                <label>Requester ID:</label>
                                <input type="text" readOnly value={ticketDetails?.RequesterID || ''} />
                            </div>
                            <div className="form-group-item">
                                <label>Assignee ID:</label>
                                <input type="text" readOnly value={ticketDetails?.AssigneeID?.Int64 || ''} />
                            </div>
                        </div>

                        <div className="form-group-grid">
                            <div className="form-group-item">
                                <label>Created At:</label>
                                <input type="text" readOnly value={ticketDetails?.CreatedAt ? new Date(ticketDetails.CreatedAt).toLocaleString() : ''} />
                            </div>
                            <div className="form-group-item">
                                <label>Updated At:</label>
                                <input type="text" readOnly value={ticketDetails?.UpdatedAt ? new Date(ticketDetails.UpdatedAt).toLocaleString() : ''} />
                            </div>
                            <div className="form-group-item">
                                <label>Closed At:</label>
                                <input type="text" readOnly value={ticketDetails?.ClosedAt?.Valid ? new Date(ticketDetails.ClosedAt.Time).toLocaleString() : ''} />
                            </div>
                        </div>
                    </div>
                )}
            </div>

            <div>
                <h3>Add Comment to My Ticket (POST /customer/tickets/{'{id}'}/comments)</h3>
                <form onSubmit={handleAddCustomerComment}>
                    <div className="form-group-item">
                        <label htmlFor="customerCommentTicketId">Ticket ID:</label>
                        <input type="text" id="customerCommentTicketId" name="customerCommentTicketId" required onInput={handleNumericInput} onKeyDown={handleNumericKeyDown} value={ticketDetails?.ID || ''} readOnly={!!ticketDetails} />
                    </div>
                    <div className="form-group-item full-width-item">
                        <label htmlFor="customerCommentBody">Comment:</label>
                        <textarea id="customerCommentBody" name="customerCommentBody" required></textarea>
                    </div>
                    <button type="submit">Add Comment</button>
                </form>
            </div>

            {ticketComments.length > 0 && (
                <div id="customerTicketCommentsSection">
                    <h3>Comments for this Ticket</h3>
                    <input type="text" id="customerCommentSearchInput" placeholder="Search comments..." style={{ width: '90%', marginBottom: '10px', marginLeft: '3%' }} onChange={(e) => setCommentSearchTerm(e.target.value)} />
                    <div style={{ maxHeight: '600px', overflowY: 'auto', border: '1px solid #0f0' }}>
                        <table id="customerTicketCommentsTable" style={{ width: '100%', marginTop: '10px' }}>
                            <thead>
                                <tr>
                                    <th style={{ border: '1px solid #0f0', padding: '8px', textAlign: 'center', position: 'sticky', top: 0, backgroundColor: '#000', zIndex: 1 }}>ID</th>
                                    <th style={{ border: '1px solid #0f0', padding: '8px', textAlign: 'center', position: 'sticky', top: 0, backgroundColor: '#000', zIndex: 1 }}>Author ID</th>
                                    <th style={{ border: '1px solid #0f0', padding: '8px', textAlign: 'left', position: 'sticky', top: 0, backgroundColor: '#000', zIndex: 1 }}>Body</th>
                                    <th style={{ border: '1px solid #0f0', padding: '8px', textAlign: 'left', position: 'sticky', top: 0, backgroundColor: '#000', zIndex: 1 }}>Created At</th>
                                </tr>
                            </thead>
                            <tbody>
                                {filteredComments.map(comment => (
                                    <tr key={comment.ID}>
                                        <td style={{ border: '1px solid #0f0', textAlign: 'center' }}>{comment.ID}</td>
                                        <td style={{ border: '1px solid #0f0', textAlign: 'center' }}>{comment.AuthorID}</td>
                                        <td style={{ border: '1px solid #0f0' }}>{comment.Body}</td>
                                        <td style={{ border: '1px solid #0f0', textAlign: 'left' }}>{new Date(comment.CreatedAt).toLocaleString()}</td>
                                    </tr>
                                ))}
                                {filteredComments.length === 0 && (
                                    <tr>
                                        <td colSpan="4" style={{ border: '1px solid #0f0' }}>No public comments found for this ticket.</td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>
            )}

            <div>
                <h3>Close My Ticket (PUT /customer/tickets/{'{id}'})</h3>
                <form onSubmit={handleCloseCustomerTicket}>
                    <div className="form-group-item">
                        <label htmlFor="customerCloseTicketId">Ticket ID:</label>
                        <input type="text" id="customerCloseTicketId" name="customerCloseTicketId" required onInput={handleNumericInput} onKeyDown={handleNumericKeyDown} />
                    </div>
                    <button type="submit">Close Ticket</button>
                </form>
            </div>
        </div>
    );
};

export default Customer;