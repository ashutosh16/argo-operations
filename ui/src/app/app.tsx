import React, { useEffect, useState} from 'react';
import { Summary } from './component/summary';
import './styles.scss';

interface SectionInfo {
    title: string;
    helpContent?: string;
}

const sectionLabel = (info: SectionInfo) => (
    <label style={{ fontSize: '12px', fontWeight: 600, color: '#6D7F8B' }}>
        {info.title}
    </label>
);

export const Extension = (props: {
    application: any;
    openFlyout: (() => void);

}) => {
    const [enabled, setEnabled] = useState(true);
    useEffect(() => {
        if (props.application.status.health === "Degraded" || props.application.status.health === "Progressing") {
            setEnabled(true);
        }
    }, [props.application.status.health]);


    const handleClick = () => {
        props.openFlyout();
    };

    return (
        <div className="application-status-panel__item">
            <div className="genai-container">
                <div className="genai-content">{sectionLabel({title: 'Intuit GenAI Failure Analysis'})}</div> {enabled ? ( <div className="genai-image-container"> <img className="genai-image" src="assets/images/genai.svg"
                                                                                                                                                                        onClick={handleClick} alt="" /> <span className="genai-footer-notes">Click to start analyzing the failed resources</span> </div> ) :
                ( <div className="genai-disabled-container">
                    <img className="genai-image" src="assets/images/genai-disabled.svg" alt=""/>
                    <span className="genai-footer-notes">AI assistant would be enabled when app health is non healthy</span> </div> )} </div>
        </div>
    );

};

export const Flyout = (props: {
    application: any;
    tree: any;
}) => {
    console.log(props.application)

    console.log(props.tree)
    return (
        <React.Fragment>
            <div>
                <Summary application={props.application} />
            </div>
        </React.Fragment>
    );
};

export const App = Extension;
