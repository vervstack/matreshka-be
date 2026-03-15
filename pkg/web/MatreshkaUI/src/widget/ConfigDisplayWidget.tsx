import {useEffect, useState} from "react";
import {
    Format,

    GetConfigRequest,

} from "@vervstack/matreshka";

import {useApi} from "@/app/hooks/api/api.ts";
import cls from "@/widget/ConfigDisplayWidget.module.css";

import LoaderWrapper from "@/segments/LoaderWrapper.tsx";
import SelectOneOf from "@/components/shared/SelectOneOf.tsx";

interface ConfigDisplayWidgetProps {
    configName: string;
}

export default function ConfigDisplayWidget({configName}: ConfigDisplayWidgetProps) {
    const {GetConfig} = useApi();

    const [format, setFormat] = useState<Format>(Format.yaml);

    const [configText, setConfigText] = useState<string>("");
    const [configType, setConfigType] = useState<string>("");

    const [loadFunc, setLoadFunc] = useState<Promise<void> | undefined>(undefined);

    const loadConfig = () => {
        const req: GetConfigRequest = {
            configName: configName,
            format: format,
        };

        const promise = GetConfig(req)
            .then(res => {
                if (res.config) {
                    setConfigText(atob(res.config.toString()));
                }

                if (res.baseInfo) {
                    res.baseInfo.configType && setConfigType(res.baseInfo.configType)

                }

            });

        setLoadFunc(promise);
    };

    useEffect(() => {
        loadConfig();
    }, [configName, format, GetConfig]);

    return (
        <div className={cls.ConfigDisplayWidgetContainer}>
            <div className={cls.Header}>
                <div className={cls.TitleSection}>
                    <div className={cls.ConfigName}>{configName}</div>
                    {configType && <div className={cls.ConfigType}>{configType}</div>}
                </div>
                <SelectOneOf
                    options={[Format.env, Format.yaml]}
                    selectedOption={format}
                    onSelected={(opt) => setFormat(opt as Format)}
                />
            </div>
            <LoaderWrapper load={loadFunc} onClickTryAgain={loadConfig}>
                <div className={cls.Content}>
                    {configText}
                </div>
            </LoaderWrapper>
        </div>
    );
}
