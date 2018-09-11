import assert from 'assert'
import {
    DidCloseTextDocumentNotification,
    DidCloseTextDocumentParams,
    DidOpenTextDocumentNotification,
    DidOpenTextDocumentParams,
    ShowInputParams,
    ShowInputRequest,
} from '../../../protocol'
import { MockConnection } from '../../../protocol/jsonrpc2/test/mockConnection'
import { URI } from '../../types/uri'
import { Window } from '../api'
import { observableValue } from '../util'
import { ExtWindows } from './windows'

describe('ExtWindows', () => {
    function create(): { extWindows: ExtWindows; mockConnection: MockConnection } {
        const mockConnection = new MockConnection()
        const extWindows = new ExtWindows(mockConnection)
        return { extWindows, mockConnection }
    }

    it('starts empty', () => {
        const { extWindows } = create()
        assert.deepStrictEqual(observableValue(extWindows), [{ isActive: true, activeComponent: null }] as Window[])
        assert.deepStrictEqual(extWindows.activeWindow, { isActive: true, activeComponent: null } as Window)
    })

    describe('component', () => {
        it('handles when a resource is opened', () => {
            const { extWindows, mockConnection } = create()
            mockConnection.recvNotification(DidOpenTextDocumentNotification.type, {
                textDocument: { uri: 'file:///a', languageId: 'l', text: 't' },
            } as DidOpenTextDocumentParams)
            const expectedWindows: Window[] = [
                { isActive: true, activeComponent: { isActive: true, resource: URI.parse('file:///a') } },
            ]
            assert.deepStrictEqual(observableValue(extWindows), expectedWindows)
            assert.deepStrictEqual(extWindows.activeWindow, expectedWindows[0])
        })

        it('handles when the open resource is closed', () => {
            const { extWindows, mockConnection } = create()
            mockConnection.recvNotification(DidOpenTextDocumentNotification.type, {
                textDocument: { uri: 'file:///a', languageId: 'l', text: 't' },
            } as DidOpenTextDocumentParams)
            mockConnection.recvNotification(DidCloseTextDocumentNotification.type, {
                textDocument: { uri: 'file:///a' },
            } as DidCloseTextDocumentParams)
            assert.deepStrictEqual(observableValue(extWindows), [{ isActive: true, activeComponent: null }] as Window[])
        })

        it('handles when a background resource is closed', () => {
            const { extWindows, mockConnection } = create()
            mockConnection.recvNotification(DidOpenTextDocumentNotification.type, {
                textDocument: { uri: 'file:///a', languageId: 'l', text: 't' },
            } as DidOpenTextDocumentParams)
            mockConnection.recvNotification(DidCloseTextDocumentNotification.type, {
                textDocument: { uri: 'file:///b' },
            } as DidCloseTextDocumentParams)
            assert.deepStrictEqual(observableValue(extWindows), [
                { isActive: true, activeComponent: { isActive: true, resource: URI.parse('file:///a') } },
            ] as Window[])
        })
    })

    describe('showInputBox', () => {
        it('sends to the client', async () => {
            const { extWindows, mockConnection } = create()
            mockConnection.mockResults.set(ShowInputRequest.type, 'c')
            const input = await extWindows.showInputBox('a', 'b')
            assert.strictEqual(input, 'c')
            assert.deepStrictEqual(mockConnection.sentMessages, [
                {
                    method: ShowInputRequest.type,
                    params: { message: 'a', defaultValue: 'b' } as ShowInputParams,
                },
            ])
        })
    })
})
