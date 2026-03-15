import {useEffect, useState} from "react";
import {useSearchParams} from "react-router-dom";
import {
    Format,

    GetConfigRequest,
    GetConfigNodeRequest,

} from "@vervstack/matreshka";

import {useApi} from "@/app/hooks/api/api.ts";
import cls from "@/widget/ConfigDisplayWidget.module.css";

import LoaderWrapper from "@/segments/LoaderWrapper.tsx";
import SelectOneOf from "@/components/shared/SelectOneOf.tsx";
import SelectDropDown from "@/components/shared/SelectDropDown.tsx";

interface ConfigDisplayWidgetProps {
    configName: string;
}

export default function ConfigDisplayWidget({configName}: ConfigDisplayWidgetProps) {
    const {GetConfig, GetConfigNodes} = useApi();

    const [searchParams, setSearchParams] = useSearchParams();
    const versionParam = searchParams.get("version");

    const [format, setFormat] = useState<Format>(Format.yaml);
    const [version, setVersion] = useState<string>(versionParam || "");
    const [versions, setVersions] = useState<string[]>([]);

    const [configText, setConfigText] = useState<string>("");
    const [configType, setConfigType] = useState<string>("");

    const [loadFunc, setLoadFunc] = useState<Promise<void> | undefined>(undefined);

    const loadConfig = () => {
        const req: GetConfigRequest = {
            configName: configName,
            format: format,
            version: version,
        };

        const promise = GetConfig(req)
            .then(res => {
                if (res.config) {
                    setConfigText(atob(res.config.toString()));
                }

                res.info?.configBase?.configType && setConfigType(res.info.configBase.configType)
            });

        setLoadFunc(promise);
    };

    useEffect(() => {
        const req: GetConfigNodeRequest = {
            configName: configName,
        };
        GetConfig(req).then(res => {
            if (res.info?.versions) {
                setVersions([...res.info.versions]);

                if (version == "") {
                    setVersion(res.info.versions[0])
                }
            }
        }).catch(err => {
            console.error("Failed to fetch versions", err);
        });
    }, [configName, GetConfigNodes]);

    useEffect(() => {
        if (version !== "") {
            setSearchParams({version}, {replace: true});
        }
    }, [version, setSearchParams]);

    useEffect(() => {
        loadConfig();
    }, [configName, format, version, GetConfig]);

    useEffect(() => {
        document.title = configName;
    }, [configName]);

    return (
        <div className={cls.ConfigDisplayWidgetContainer}>
            <div className={cls.Header}>
                <div className={cls.TitleSection}>
                    <div className={cls.ConfigName}>{configName}</div>
                    {configType && <div className={cls.ConfigType}>{configType}</div>}
                </div>
                <div className={cls.Actions}>
                    <div className={cls.VersionSelect}>
                        <SelectDropDown
                            label="Version"
                            options={versions}
                            selectedOption={version}
                            onSelected={setVersion}
                        />
                    </div>
                    <div className={cls.FormatSelect}>
                        <SelectOneOf
                            options={[Format.env, Format.yaml]}
                            selectedOption={format}
                            onSelected={(opt) => setFormat(opt as Format)}
                        />
                    </div>
                </div>
            </div>

            <LoaderWrapper load={loadFunc} onClickTryAgain={loadConfig}>
                <div className={cls.Content}>
                    {configText}
                </div>
            </LoaderWrapper>
        </div>
    );
}
