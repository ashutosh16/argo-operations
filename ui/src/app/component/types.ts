import {models} from 'argo-ui';

export type ArgoOperationResult = {
    apiVersion: string;
    kind: string;
    metadata: models.ObjectMeta;
    spec?: ArgoOperationSpec;
    status?: {
        startedAt: string;
        finishedAt: string;
        phase: string;
        message: string;
        OperationResults: OperationResult[];
    };
};

export type ArgoOperationSpec = {
    name?: string;
    args?: { name: string; value?: string }[];
    feedback?: boolean;
    help: {
        stackOverflow?: string[];
        slack: {
            primary: string;
            secondary: string;
        };
    };
};


export type OperationResult = {
    name: string;
    finishedAt: string;
    startedAt: string;
    summary: { [key: string]: string };
    feedback?: { [key: string]: string };
};