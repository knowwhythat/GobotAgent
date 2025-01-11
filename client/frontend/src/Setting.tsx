import {Button, Flex, Input, message, Modal, Space, Typography} from 'antd';
import React, {useEffect, useState} from "react";
import {FolderOpenOutlined} from "@ant-design/icons";
import {GetSettingData, OpenFileDialog, SetSettingData} from "../wailsjs/go/main/App";
import {models} from "../wailsjs/go/models.ts";
import SettingData = models.SettingData;

interface SettingProps {
    isOpen: boolean;
    onClose: () => void;
}

const Setting: React.FC<SettingProps> = ({isOpen, onClose}) => {
    const [settingData, setSettingData] = useState<SettingData>({
        browserPath: "",
        apiBaseUrl: "",
        apiKey: "",
        modelName: ""
    })

    useEffect(() => {
        const fetchData = async () => {
            try {
                const data = await GetSettingData();
                setSettingData(data);
            } catch (error) {
                console.error("Failed to fetch setting data:", error);
            }
        };
        fetchData().finally();
    }, [])

    const handleValueChange = (name: string, value: string) => {
        setSettingData(prevState => {
            return {
                ...prevState,
                [name]: value
            }
        })
    }

    const handleSubmit = async () => {
        await SetSettingData(settingData)
        onClose()
        message.success("设置保存成功")
    }

    const handleOpenFileDialog = async () => {
        const path = await OpenFileDialog([{DisplayName: "浏览器路径", Pattern: "*.exe"}])
        if (path) {
            handleValueChange("browserPath", path)
        }
    }

    return (
        <div>
            <Modal
                width={800}
                title="设置"
                open={isOpen}
                onOk={handleSubmit}
                onCancel={onClose}
            >
                <Flex vertical={true} gap={16} className={"mt-3"}>
                    <div>
                        <Typography.Text className={"w-32"}>浏览器路径</Typography.Text>
                        <Space.Compact style={{width: '100%'}}>
                            <Input name="browserPath" value={settingData.browserPath} onChange={(e) => {
                                handleValueChange(e.target.name, e.target.value)
                            }} placeholder={"浏览器路径"}/>
                            <Button icon={<FolderOpenOutlined/>} onClick={handleOpenFileDialog}></Button>
                        </Space.Compact>
                    </div>
                    <div>
                        <Typography.Text className={"w-32"}>API_BASE_URL</Typography.Text>
                        <Space.Compact style={{width: '100%'}}>
                            <Input name="apiBaseUrl" value={settingData.apiBaseUrl} onChange={(e) => {
                                handleValueChange(e.target.name, e.target.value)
                            }} placeholder={"API_BASE_URL"}/>
                        </Space.Compact>
                    </div>
                    <div>
                        <Typography.Text className={"w-32"}>API_KEY</Typography.Text>
                        <Space.Compact style={{width: '100%'}}>
                            <Input name="apiKey" value={settingData.apiKey} onChange={(e) => {
                                handleValueChange(e.target.name, e.target.value)
                            }} placeholder={"API_KEY"}/>
                        </Space.Compact>
                    </div>
                    <div>
                        <Typography.Text className={"w-32"}>模型名称</Typography.Text>
                        <Space.Compact style={{width: '100%'}}>
                            <Input name="modelName" value={settingData.modelName} onChange={(e) => {
                                handleValueChange(e.target.name, e.target.value)
                            }} placeholder={"模型名称"}/>
                        </Space.Compact>
                    </div>
                </Flex>
            </Modal>
        </div>
    )
}

export default Setting;