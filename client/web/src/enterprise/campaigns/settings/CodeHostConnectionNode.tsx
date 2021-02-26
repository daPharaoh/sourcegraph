import React, { useCallback, useEffect, useState } from 'react'
import * as H from 'history'
import CheckboxBlankCircleOutlineIcon from 'mdi-react/CheckboxBlankCircleOutlineIcon'
import CheckCircleOutlineIcon from 'mdi-react/CheckCircleOutlineIcon'
import { defaultExternalServices } from '../../../components/externalServices/externalServices'
import {
    CampaignsCodeHostFields,
    CampaignsCredentialFields,
    ExternalServiceKind,
    Scalars,
} from '../../../graphql-operations'
import { AddCredentialModal } from './AddCredentialModal'
import { RemoveCredentialModal } from './RemoveCredentialModal'
import { Subject } from 'rxjs'
import Dialog from '@reach/dialog'
import ContentCopyIcon from 'mdi-react/ContentCopyIcon'
import copy from 'copy-to-clipboard'

export interface CodeHostConnectionNodeProps {
    node: CampaignsCodeHostFields
    userID: Scalars['ID']
    history: H.History
    updateList: Subject<void>
}

type OpenModal = 'add' | 'view' | 'delete'

export const CodeHostConnectionNode: React.FunctionComponent<CodeHostConnectionNodeProps> = ({
    node,
    userID,
    history,
    updateList,
}) => {
    const Icon = defaultExternalServices[node.externalServiceKind].icon

    const [openModal, setOpenModal] = useState<OpenModal | undefined>()
    const onClickAdd = useCallback(() => {
        setOpenModal('add')
    }, [])
    const onClickRemove = useCallback<React.MouseEventHandler>(event => {
        event.preventDefault()
        setOpenModal('delete')
    }, [])
    const onClickView = useCallback<React.MouseEventHandler>(event => {
        event.preventDefault()
        setOpenModal('view')
    }, [])
    const closeModal = useCallback(() => {
        setOpenModal(undefined)
    }, [])
    const afterAction = useCallback(() => {
        setOpenModal(undefined)
        updateList.next()
    }, [updateList])

    const isEnabled = node.credential !== null

    return (
        <>
            <li className="list-group-item p-3 test-code-host-connection-node">
                <div className="d-flex justify-content-between align-items-center mb-0">
                    <h3 className="mb-0">
                        {isEnabled && (
                            <CheckCircleOutlineIcon
                                className="text-success icon-inline test-code-host-connection-node-enabled"
                                data-tooltip="Connected"
                            />
                        )}
                        {!isEnabled && (
                            <CheckboxBlankCircleOutlineIcon
                                className="text-danger icon-inline test-code-host-connection-node-disabled"
                                data-tooltip="No token set"
                            />
                        )}
                        <Icon className="icon-inline mx-2" /> {node.externalServiceURL}
                    </h3>
                    <div className="mb-0">
                        {isEnabled && (
                            <>
                                <a
                                    href=""
                                    className="btn btn-link text-danger test-code-host-connection-node-btn-remove"
                                    onClick={onClickRemove}
                                >
                                    Remove
                                </a>
                                {node.requiresSSH && (
                                    <button type="button" onClick={onClickView} className="btn btn-secondary ml-2">
                                        View public key
                                    </button>
                                )}
                            </>
                        )}
                        {!isEnabled && (
                            <button
                                type="button"
                                className="btn btn-success test-code-host-connection-node-btn-add"
                                onClick={onClickAdd}
                            >
                                Add token
                            </button>
                        )}
                    </div>
                </div>
            </li>
            {openModal === 'delete' && (
                <RemoveCredentialModal
                    onCancel={closeModal}
                    afterDelete={afterAction}
                    history={history}
                    codeHost={node}
                    credential={node.credential!}
                />
            )}
            {openModal === 'view' && (
                <ViewCredentialModal onClose={closeModal} codeHost={node} credential={node.credential!} />
            )}
            {openModal === 'add' && (
                <AddCredentialModal
                    onCancel={closeModal}
                    afterCreate={afterAction}
                    history={history}
                    userID={userID}
                    externalServiceKind={node.externalServiceKind}
                    externalServiceURL={node.externalServiceURL}
                    requiresSSH={node.requiresSSH}
                />
            )}
        </>
    )
}

interface ViewCredentialModalProps {
    codeHost: CampaignsCodeHostFields
    credential: CampaignsCredentialFields

    onClose: () => void
}

export const ViewCredentialModal: React.FunctionComponent<ViewCredentialModalProps> = ({
    credential,
    codeHost,
    onClose,
}) => {
    const labelId = 'viewCredential'
    return (
        <Dialog
            className="modal-body modal-body--top-third p-4 rounded border"
            onDismiss={onClose}
            aria-labelledby={labelId}
        >
            <div className="test-remove-credential-modal">
                <h3 id={labelId}>
                    Campaigns credentials: {defaultExternalServices[codeHost.externalServiceKind].defaultDisplayName}
                </h3>
                <p>
                    <strong>{codeHost.externalServiceURL}</strong>
                </p>

                <h4>Personal access token</h4>
                <p>
                    <i>PATs cannot be viewed after entering.</i>
                </p>

                <hr className="mb-3" />

                <CodeHostSSHPublicKey
                    externalServiceKind={codeHost.externalServiceKind}
                    sshPublicKey={credential.sshPublicKey!}
                />

                <div className="d-flex justify-content-end pt-5">
                    <button type="button" className="btn btn-outline-secondary" onClick={onClose}>
                        Close
                    </button>
                </div>
            </div>
        </Dialog>
    )
}

export interface CodeHostSSHPublicKeyProps {
    externalServiceKind: ExternalServiceKind
    sshPublicKey: string
    label?: string
    showInstructionsLink?: boolean
    showCopyButton?: boolean
}

const configInstructionLinks: Record<ExternalServiceKind, string> = {
    [ExternalServiceKind.GITHUB]:
        'https://docs.github.com/en/github/authenticating-to-github/adding-a-new-ssh-key-to-your-github-account',
    [ExternalServiceKind.GITLAB]: 'https://docs.gitlab.com/ee/ssh/#add-an-ssh-key-to-your-gitlab-account',
    [ExternalServiceKind.BITBUCKETSERVER]:
        'https://confluence.atlassian.com/bitbucketserver/ssh-user-keys-for-personal-use-776639793.html',
    [ExternalServiceKind.AWSCODECOMMIT]: 'unsupported',
    [ExternalServiceKind.BITBUCKETCLOUD]: 'unsupported',
    [ExternalServiceKind.GITOLITE]: 'unsupported',
    [ExternalServiceKind.OTHER]: 'unsupported',
    [ExternalServiceKind.PERFORCE]: 'unsupported',
    [ExternalServiceKind.PHABRICATOR]: 'unsupported',
}

export const CodeHostSSHPublicKey: React.FunctionComponent<CodeHostSSHPublicKeyProps> = ({
    externalServiceKind,
    sshPublicKey,
    showInstructionsLink = true,
    showCopyButton = true,
    label = 'Public SSH key',
}) => {
    const [copied, setCopied] = useState<boolean>(false)
    const onCopy = useCallback(() => {
        copy(sshPublicKey)
        setCopied(true)
    }, [sshPublicKey])
    useEffect(() => {
        if (copied) {
            const timer = setTimeout(() => {
                setCopied(false)
            }, 1500)
            return () => clearTimeout(timer)
        }
        return () => undefined
    }, [copied])
    return (
        <>
            <div className="d-flex justify-content-between align-items-end mb-2">
                <h4>{label}</h4>
                {showCopyButton && (
                    <button type="button" className="btn btn-secondary" onClick={onCopy}>
                        <ContentCopyIcon className="icon-inline" />
                        {copied ? 'Copied!' : 'Copy'}
                    </button>
                )}
            </div>
            <textarea className="form-control text-monospace mb-3" rows={5} value={sshPublicKey} />
            {showInstructionsLink && (
                <p>
                    <a href={configInstructionLinks[externalServiceKind]} target="_blank" rel="noopener">
                        Configuration instructions
                    </a>
                </p>
            )}
        </>
    )
}
