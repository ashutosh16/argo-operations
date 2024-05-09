import React, { useState } from 'react';
import axios from 'axios';
import './feedback.scss';

export const FeedbackComponent = (props) => {
    const [showResponseBox, setShowResponseBox] = useState(false);
    const [feedbackText, setFeedbackText] = useState('');

    const handleTumbdownClick = () => {
        setShowResponseBox(true);
    };

    const handleTextChange = (event) => {
        setFeedbackText(event.target.value);
    };

    const handleSubmit = () => {
        axios
            .post('api/feedback', { text: feedbackText })
            .then(() => {
                setShowResponseBox(false);
                setFeedbackText('');
            })
            .catch((error) => {
                console.error(error);
            });
    };

    return (
        <div className="feedback">
            <div className={`votes ${showResponseBox ? 'disabled' : ''}`}>
                <a className="upvote" onClick={handleTumbdownClick}>
                    <i className="fas fa-thumbs-up"></i>
                </a>
                <a className="downvote" onClick={handleTumbdownClick}>
                    <i className="fas fa-thumbs-down"></i>
                </a>
            </div>
            {showResponseBox && (
                <div className="modal-content">
                    <div className="modal-header">
                        Would you like to tell us more?
                        <i className="fas fa-times" onClick={() => setShowResponseBox(false)}></i>
                    </div>
                    <div className="modal-body">
                        <textarea id="modalCommentArea" className="modal-comment-area" value={feedbackText} onChange={handleTextChange} placeholder="Enter your feedback here..." />
                        <button className="modal-submit-button" onClick={handleSubmit}>
                            <span>Send</span>
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
};