import * as React from 'react';
import * as models from "argocd/ui/src/app/shared/models"
import './resourcecard.scss';

interface ResourceCardsProps{
app: models.Application;
}
export const ResourceCards = (props: ResourceCardsProps) => {
    const resNode: models.ResourceStatus[] = props.app.status.resources.filter(res => res.health?.status == models.HealthStatuses.Degraded || res.health?.status == models.HealthStatuses.Progressing || res.health?.status == models.HealthStatuses.Unknown);
    return (
        <>
            {resNode.map(res => (
                <div
                    className='applications-tiles argo-table-list argo-table-list--clickable row small-up-1 medium-up-2 large-up-3 xxxlarge-up-4'>
                    <div
                        className={`argo-table-list__row applications-list__entry applications-list__entry--health-${res.health?.status}`}>
                        <div
                            className={`columns small-12 applications-list__info  applications-tiles__item`}>
                            <div className='row '>
                                Name: {res?.name}
                            </div>
                            <div className='row '>
                                Health Details: {res?.health?.message}
                            </div>
                            <div className='row '>
                                status: {res?.health?.status}
                            </div>
                        </div>
                    </div>
                </div>
            ))}
        </>
    );
}