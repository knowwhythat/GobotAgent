import type {BubbleProps} from '@ant-design/x';
import {Bubble, Conversations, ConversationsProps, Prompts, Sender, ThoughtChain, XProvider,} from '@ant-design/x';
import {Button, Divider, Flex, GetProp, message, Space, Spin, Splitter, Typography} from 'antd';
import {useState} from 'react';

import {
    ChromeOutlined,
    DeleteOutlined,
    EditOutlined,
    FileAddOutlined,
    OpenAIOutlined,
    SettingOutlined,
    UserOutlined
} from '@ant-design/icons';
import {nanoid} from "nanoid";
import {Greet} from "../wailsjs/go/main/App";
import {EventsOff, EventsOn} from "../wailsjs/runtime";
import {AgentResponse} from "./types.ts";
import markdownit from 'markdown-it';
import Setting from "./Setting.tsx";

const md = markdownit({html: true, breaks: true});

const renderMarkdown: BubbleProps['messageRender'] = (content) => (
    <Typography>
        <div dangerouslySetInnerHTML={{__html: md.render(content)}}/>
    </Typography>
);

const roles: GetProp<typeof Bubble.List, 'roles'> = {
    ai: {
        placement: 'start',
        avatar: {icon: <OpenAIOutlined/>, style: {background: '#2d539c'}},
        typing: {step: 5, interval: 20},
        style: {
            maxWidth: 600,
        },
    },
    user: {
        placement: 'end',
        avatar: {icon: <UserOutlined/>, style: {background: '#87d068'}},
    },
};


const ChatApp = () => {
    const [value, setValue] = useState('');
    const [loading, setLoading] = useState(false);
    const [bubbleListItems, setBubbleListItems] = useState<GetProp<typeof Bubble.List, 'items'>>([{
        key: nanoid(8),
        role: 'ai',
        content: '你好，请在下方输入你的需求',

    },]);
    const [conversationListItems, setConversationListItems] = useState<GetProp<typeof Conversations, 'items'>>([
        {
            key: nanoid(8),
            label: '会话 - 1',
            icon: <ChromeOutlined/>,
        },
    ]);
    const [thoughtChainItems, setThoughtChainItems] = useState<GetProp<typeof ThoughtChain, 'items'>>([]);
    const [isModalOpen, setIsModalOpen] = useState(false);

    async function handleSubmit() {
        setLoading(true);
        // 先将用户的输入更新到 bubbleListItems
        const newUserItem = {key: nanoid(8), role: 'user', content: value};
        setValue('');
        setBubbleListItems(prevItems => [...prevItems, newUserItem]);
        try {
            await Greet(value);
            await new Promise((resolve) => {
                EventsOn("agent_response", (data: string) => {
                    const resp = JSON.parse(data) as AgentResponse
                    if (resp.done) {
                        setThoughtChainItems(prevSItems => [...prevSItems, {
                            key: nanoid(8),
                            icon: <ChromeOutlined/>,
                            title: `完成`,
                        }]);
                        setBubbleListItems(prevItems => [...prevItems, {
                            key: nanoid(8),
                            role: 'ai',
                            content: resp.result,
                            messageRender: renderMarkdown
                        }]);
                        EventsOff("agent_response")
                        resolve(true)
                    } else {
                        setThoughtChainItems(prevSItems => [...prevSItems, {
                            key: nanoid(8),
                            icon: <ChromeOutlined/>,
                            title: `步骤${prevSItems.length + 1}`,
                            content: resp.result,
                        }]);
                    }
                })
            });
        } catch (e: any) {
            message.error(e);
        } finally {
            setLoading(false);
        }
    }

    function handleSetting() {
        setIsModalOpen(true)
    }

    function createConversation() {
        const newConversationItem = {
            key: nanoid(8),
            label: `会话 - ${conversationListItems.length + 1}`,
            icon: <ChromeOutlined/>
        };
        setConversationListItems(prevItems => [...prevItems, newConversationItem])
    }

    const menuConfig: ConversationsProps['menu'] = (conversation) => ({
        items: [
            {
                label: '重命名',
                key: 'rename',
                icon: <EditOutlined/>,
            },
            {
                label: '删除',
                key: 'delete',
                icon: <DeleteOutlined/>,
                danger: true,
            },
        ],
        onClick: async (menuInfo) => {
            await message.success(`Click ${conversation.key} - ${menuInfo.key}`);
        },
    });

    return (
        <>
            <XProvider>
                <Flex className={'w-full h-full py-4 px-2'} gap={12}>
                    <Splitter className={"h-full"}>
                        <Splitter.Panel defaultSize="80%" min="50%">
                            <Flex className={'w-full h-full py-4 px-2'}>
                                <Flex vertical={true} className={"h-full"}>
                                    <Conversations
                                        className={"w-40 flex-1"}
                                        defaultActiveKey="1"
                                        menu={menuConfig}
                                        items={conversationListItems}
                                    />
                                    <Button type="primary" icon={<SettingOutlined/>}
                                            onClick={handleSetting}>设置</Button>
                                    <Setting isOpen={isModalOpen} onClose={() => setIsModalOpen(false)}/>
                                </Flex>
                                <Divider type="vertical" className={"h-full"}/>
                                <Flex vertical className={"h-full flex-1"} gap={8}>
                                    <Bubble.List
                                        className={"flex-1"}
                                        autoScroll={true}
                                        roles={roles}
                                        items={bubbleListItems}
                                    />
                                    <Prompts
                                        items={[
                                            {
                                                key: nanoid(8),
                                                icon: <FileAddOutlined className={"text-fuchsia-500"}/>,
                                                label: '创建新会话'
                                            }
                                        ]}
                                        onItemClick={createConversation}
                                    />
                                    <Sender
                                        submitType="enter"
                                        value={value}
                                        loading={loading}
                                        onChange={setValue}
                                        onSubmit={handleSubmit}
                                        actions={(_, info) => {
                                            const {SendButton, LoadingButton, ClearButton} = info.components;

                                            return (
                                                <Space size="small">
                                                    <Typography.Text type="secondary">
                                                        <small>`Shift + Enter` to newline</small>
                                                    </Typography.Text>
                                                    <ClearButton/>
                                                    {loading ? (
                                                        <LoadingButton type="default" icon={<Spin size="small"/>}
                                                                       disabled/>
                                                    ) : (
                                                        <SendButton type="primary" icon={<OpenAIOutlined/>}
                                                                    disabled={false}/>
                                                    )}
                                                </Space>
                                            );
                                        }}
                                    />
                                </Flex>
                            </Flex>
                        </Splitter.Panel>
                        <Splitter.Panel>
                            <Flex vertical={true} className={"h-full"}>
                                <p className={"w-full text-center select-none"}>执行步骤</p>
                                <ThoughtChain
                                    className={"flex-1 overflow-auto"}
                                    items={thoughtChainItems}
                                />
                            </Flex>
                        </Splitter.Panel>
                    </Splitter>
                </Flex>
            </XProvider>
        </>
    );
};
export default ChatApp;