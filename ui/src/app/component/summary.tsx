import * as React from "react";
import { useState, useEffect } from "react";
import { FeedbackComponent } from "./feedback/feedback";
import './summary.scss';
import {ResourceCards} from "./resourcecard/resourcecard";



const Target_Status = ['Degraded', 'Progressing', "Healthy"];
const GenAI_Res_Name = "genai-analysis"
export const Summary: any = ({
                                 application,
                             }) => {
    const [analysisOperationResultSummary, setAnalysisOperationResultSummary] = useState<any>(null);
    const [initiateAction, setInitiateAction] = useState(false);

    const applicationName = application?.metadata?.name || "";
    const applicationNamespace = application?.metadata?.namespace || "";
    const destNamespace = application.spec.destination.namespace || "";
   // initiation will be called when gen-ai object has live longer than 5 mins


        const fetchExistingGenAiResource = async () => {
            const url = `/api/v1/applications/${applicationName}/resource?name=${GenAI_Res_Name}&appNamespace=${applicationNamespace}&resourceName=${GenAI_Res_Name}&namespace=${destNamespace}&version=v1alpha1&kind=ArgoOperationResult&group=argoproj.io`;
            try {
                const response = await fetch(url);
                if (response.status  === 404) {
                    setInitiateAction(false)
                }
                const data = await response.json();
                const json = typeof data.manifest === "string" ? JSON.parse(data.manifest) : data.manifest;
                const finishedTime = new Date(json.status?.finishedAt.replace('Z', '')).getTime(); //parse the date and remove 'Z'
                const currentTime = new Date().getTime();
                const timeDifference = currentTime - finishedTime;
                const minutesPassed = timeDifference / (1000 * 60);

                if (minutesPassed > 5) {
                    setAnalysisOperationResultSummary(null);
                    setInitiateAction(true)
                }else{
                    setAnalysisOperationResultSummary(json?.status?.operationResults?.find((result) => result?.name === "genai-analysis"));
                    setInitiateAction(false)
                }

            } catch (err) {
                setInitiateAction(true)
                console.error("Error fetching analysisOperationResultSummary:", err);
            }
            finally {

            }
        };

    useEffect(() => {
        fetchExistingGenAiResource()
        console.log("initiateAction" +initiateAction)
        if (!initiateAction) return;

        const rollouts = application.status.resources.filter(resource => {
            return resource.kind === "Rollout" && Target_Status.includes(resource.health?.status);
        });
        const rollout = rollouts.length > 1 ? rollouts : rollouts[0];

        const url = `/api/v1/applications/${applicationName}/resource/actions?appNamespace=${applicationNamespace}&namespace=${destNamespace}&resourceName=${rollout?.name}&version=v1alpha1&kind=Rollout&group=argoproj.io`;

        const action = async () => {
            try {
                await fetch(url, {
                    method: 'POST',
                    headers: {
                        Accept: '*/*',
                        'Accept-Language': 'en-US,en;q=0.9',
                        Connection: 'keep-alive',
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify("genai-analysis"),
                });

                // Delay the execution of getAnalysisOperationResultSummary
                const delay = 30000;
                const timeoutId = setTimeout(async () => {
                    try {
                        const response = await fetch(`/api/v1/applications/${applicationName}/resource?name=${GenAI_Res_Name}&appNamespace=${applicationNamespace}&resourceName=${GenAI_Res_Name}&namespace=${destNamespace}&version=v1alpha1&kind=ArgoOperationResult&group=argoproj.io`);
                        if (response.status > 499) {
                            throw new Error("No metrics");
                        }
                        const data = await response.json();
                        const json = typeof data.manifest === "string" ? JSON.parse(data.manifest) : data.manifest;
                        setAnalysisOperationResultSummary(json?.status?.operationResults?.find((result) => result?.name === "genai-analysis"));
                    } catch (err) {
                        console.error("Error fetching analysisOperationResultSummary:", err);
                    }
                    finally {
                        setInitiateAction(false)
                    }
                }, delay);

                // Cleanup function to clear the timeout if the component unmounts or it's time to trigger the next iteration of the API call
                return () => clearTimeout(timeoutId);

            } catch (err) {
                console.error('Error calling initiate action:', err);
            }
        };

        action();

    }, [initiateAction]);

console.log('analysisOperationResultSummary+'+analysisOperationResultSummary)

    return (
        <React.Fragment>
            {analysisOperationResultSummary != null ? (
                <div>
                    <div className="summary__header">
                        <h1>Summary</h1>
                    </div>
                    <button   className={`argo-button argo-button--base'${initiateAction ? "summary__btn-enable" : "summary__btn-disable"}`} onClick={fetchExistingGenAiResource}>
                        Summarize again
                    </button>
                    <div className="summary__feedback">
                        <FeedbackComponent />
                    </div>
                    <div className="summary__content">
                        <div className="summary__row">
                            <div className="summary__row__label">
                                Overview
                            </div>
                            <div className="summary__row__value">
                                {analysisOperationResultSummary.summary.overallSummary}
                            </div>
                        </div>
                        <div style={{paddingTop: '20px' }} className="summary__row">
                            <div className="summary__row__label">
                                GenAI Recommendation
                            </div>
                            <div className="summary__row__value">
                                {analysisOperationResultSummary.summary.aiRecommendation}
                            </div>
                        </div>
                        <p className="warning">
                            <i className="fa fa-exclamation-triangle" /> Summary may display inaccurate info, including about deployment, so double-check its responses
                        </p>
                        <div style={{paddingTop: "50px"}} className="summary__row">
                            <ResourceCards app={application} />
                        </div>
                    </div>
                </div>
            ) : (
                <> <i className="fa-solid fa-spinner"></i>
                    <span>Fetching the Data</span>
                </>

            )}
        </React.Fragment>
    );
};
