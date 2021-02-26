import React, { useCallback, useState } from 'react'
import * as H from 'history'
import Dialog from '@reach/dialog'
import { LoadingSpinner } from '@sourcegraph/react-loading-spinner'
import { asError, isErrorLike } from '../../../../../shared/src/util/errors'
import { ErrorAlert } from '../../../components/alerts'
import { deleteCampaignsCredential } from './backend'
import { CampaignsCodeHostFields, CampaignsCredentialFields } from '../../../graphql-operations'
import { defaultExternalServices } from '../../../components/externalServices/externalServices'
import { CodeHostSSHPublicKey } from './CodeHostConnectionNode'

export interface RemoveCredentialModalProps {
    codeHost: CampaignsCodeHostFields
    credential: CampaignsCredentialFields

    onCancel: () => void
    afterDelete: () => void

    history: H.History
}

export const RemoveCredentialModal: React.FunctionComponent<RemoveCredentialModalProps> = ({
    codeHost,
    credential,
    onCancel,
    afterDelete,
    history,
}) => {
    const labelId = 'removeCredential'
    const [isLoading, setIsLoading] = useState<boolean | Error>(false)
    const onDelete = useCallback<React.MouseEventHandler>(async () => {
        setIsLoading(true)
        try {
            await deleteCampaignsCredential(credential.id)
            afterDelete()
        } catch (error) {
            setIsLoading(asError(error))
        }
    }, [afterDelete, credential.id])
    return (
        <Dialog
            className="modal-body modal-body--top-third p-4 rounded border"
            onDismiss={onCancel}
            aria-labelledby={labelId}
        >
            <div className="test-remove-credential-modal">
                <h3>
                    Campaigns credentials: {defaultExternalServices[codeHost.externalServiceKind].defaultDisplayName}
                </h3>
                <p>
                    <strong>{codeHost.externalServiceURL}</strong>
                </p>
                <h3 className="text-danger" id={labelId}>
                    Removing credentials is irreversible
                </h3>

                {isErrorLike(isLoading) && <ErrorAlert error={isLoading} history={history} />}

                <p>
                    To create changesets on this code host after removing credentials, you will need to repeat the 'Add
                    credentials' process.
                </p>

                {codeHost.requiresSSH && (
                    <CodeHostSSHPublicKey
                        externalServiceKind={codeHost.externalServiceKind}
                        sshPublicKey={credential.sshPublicKey!}
                        showInstructionsLink={false}
                        showCopyButton={false}
                        label="Public key to remove"
                    />
                )}

                <div className="d-flex justify-content-end pt-5">
                    <button
                        type="button"
                        disabled={isLoading === true}
                        className="btn btn-outline-secondary mr-2"
                        onClick={onCancel}
                    >
                        Cancel
                    </button>
                    <button
                        type="button"
                        disabled={isLoading === true}
                        className="btn btn-danger test-remove-credential-modal-submit"
                        onClick={onDelete}
                    >
                        {isLoading === true && <LoadingSpinner className="icon-inline" />}
                        Remove credentials
                    </button>
                </div>
            </div>
        </Dialog>
    )
}
