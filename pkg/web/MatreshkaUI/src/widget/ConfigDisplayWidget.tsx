import {useEffect, useState} from "react";
import {useSearchParams} from "react-router-dom";
import {
    Format,

    GetConfigRequest,
    GetConfigNodeRequest,
    SaveConfigRequest,
} from "@vervstack/matreshka";

import {useApi} from "@/app/hooks/api/api.ts";
import cls from "@/widget/ConfigDisplayWidget.module.css";
import Pencil from "@/assets/icons/Pencil.svg";
import X from "@/assets/icons/X.svg";

import LoaderWrapper from "@/segments/LoaderWrapper.tsx";
import SelectOneOf from "@/components/shared/SelectOneOf.tsx";
import SelectDropDown from "@/components/shared/SelectDropDown.tsx";


interface ConfigDisplayWidgetProps {
    configName: string;
}

export default function ConfigDisplayWidget({configName}: ConfigDisplayWidgetProps) {
    const {GetConfig, GetConfigNodes, SaveConfig} = useApi();

    const [searchParams, setSearchParams] = useSearchParams();
    const versionParam = searchParams.get("version");

    const [format, setFormat] = useState<Format>(Format.yaml);
    const [version, setVersion] = useState<string>(versionParam || "");
    const [versions, setVersions] = useState<string[]>([]);

    const [configText, setConfigText] = useState<string>("");
    const [editedConfigText, setEditedConfigText] = useState<string>("");


    const [isEditing, setIsEditing] = useState<boolean>(false);
    const [configType, setConfigType] = useState<string>("");

    const lineCount = (isEditing ? editedConfigText : configText).split("\n").length;
    const editorHeight = `${lineCount * 1.5}em`;

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
                    const text = atob(res.config.toString());
                    setConfigText(text);
                    setEditedConfigText(text);
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

                if (!res.info.versions.includes(version)) {
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
        if (version != "") {
            document.title += `@${version}`;
        }
    }, [configName, version]);

    const isDirty = configText !== editedConfigText;

    const onSave = () => {
        const req: SaveConfigRequest = {
            configName: configName,
            format: format,
            version: version,
            config: new TextEncoder().encode(editedConfigText),
        };

        SaveConfig(req).then(() => {
            setConfigText(editedConfigText);
            setIsEditing(false);
        }).catch(err => {
            console.error("Failed to save config", err);
        });
    };

    const onRollback = () => {
        setEditedConfigText(configText);
        setIsEditing(false);
    };

    return (
        <div className={cls.ConfigDisplayWidgetContainer}>
            <div className={cls.Header}>
                <div className={cls.TitleSection}>
                    <div className={cls.ConfigName}>{configName}</div>
                    {configType && <div className={cls.ConfigType}>{configType}</div>}
                </div>
                <div className={cls.Actions}>
                    {isDirty && (
                        <div className={cls.EditActions}>
                            <button className={cls.RollbackButton} onClick={onRollback}>rollback</button>
                            <button className={cls.SaveButton} onClick={onSave}>save</button>
                        </div>
                    )}
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
                    <div className={cls.EditIcon}
                         onClick={() => setIsEditing(!isEditing)}>
                        {isEditing ? (<img src={X} alt="Stop editing"/>)
                            :
                            (<img src={Pencil} alt="Edit"/>)}
                    </div>

                    {isEditing ? (
                        <textarea
                            className={cls.Editor}
                            value={editedConfigText}
                            onChange={(e) => setEditedConfigText(e.target.value)}
                            style={{height: editorHeight}}
                            autoFocus
                        />
                    ) : (
                        <pre className={cls.Pre} style={{height: editorHeight}}>{configText}</pre>
                    )}
                </div>
            </LoaderWrapper>
        </div>
    );
}
