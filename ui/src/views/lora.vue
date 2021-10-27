<template>
    <div class="main">
        <div class="model">
            <v-input class="m-left" type="warn" v-model:value="apiKey"></v-input>
            <v-btn class="m-right" type="warn" size="large" @click="saveKey">添加 API Key</v-btn>
        </div>
        <div class="model">
            <v-input class="m-left" type="warn" v-model:value="bridgeId"></v-input>
            <v-btn class="m-right" type="warn" size="large" @click="saveBridge">添加 ApplicationID</v-btn>
        </div>
        <div class="model">
            <v-input class="m-left" type="warn" v-model:value="lockId"></v-input>
            <v-btn class="m-right" type="warn" size="large" @click="saveLock">添加 DeviceID</v-btn>
        </div>
        <div class="model">
            <template v-for="(key, index) in data.apis" :key="index">
                <v-radio type="success" :value="key" v-model:checkedValue="apiKeySelect">{{ key }}</v-radio>
            </template>
        </div>
        <div class="model">
            <template v-for="(key, index) in data.bridges" :key="index">
                <v-radio type="pink" :value="key" v-model:checkedValue="bridgeIdSelect">{{ key }}</v-radio>
            </template>
        </div>
        <div class="model">
            <template v-for="(key, index) in data.locks" :key="index">
                <v-radio type="purple" :value="key" v-model:checkedValue="lockIdSelect">{{ key }}</v-radio>
            </template>
        </div>
        <div class="model">
            <v-input class="m-left" type="success" v-model:value="message"></v-input>
            <v-btn class="m-right" type="success" size="large" @click="sendMessage">发送消息</v-btn>
        </div>
        <div class="model">
            <v-btn size="large" type="purple" @click="fnLock">开锁</v-btn>
            <v-btn size="large" type="purple" @click="fnUnlock">上锁</v-btn>
            <v-btn size="large" type="purple" @click="fnBattery">查询电量</v-btn>
            <v-btn size="large" type="purple" @click="fnState">查询状态</v-btn>
            <!-- <v-btn size="large" type="purple" @click="fnMessage">准备获取信息</v-btn> -->
            <v-btn size="large" type="purple" @click="connWS">连接WebSocket</v-btn>
            <v-btn size="large" @click="clearLogs">清理日志</v-btn>
        </div>
        <div class="logs" ref="logs"></div>
    </div>
</template>

<script lang="ts">
import { defineComponent, ref, reactive } from 'vue';
import axios from 'axios';
export default defineComponent({
    setup(props, context) {
        axios.defaults.headers.post['Content-Type'] = 'application/json';
        const apiKey = ref<string>('');
        const bridgeId = ref<string>('');
        const lockId = ref<string>('');
        const getTime = (): string => {
            var d = new Date(),
                h = d.getHours(),
                m = d.getMinutes(),
                s = d.getSeconds();
            return `${h < 10 ? '0' + h : h}:${m < 10 ? '0' + m : m}:${s < 10 ? '0' + s : s}`;
        };
        const data = reactive({
            apis: Array<string>(),
            bridges: Array<string>(),
            locks: Array<string>()
        });
        if (localStorage.getItem('apis22')) data.apis = JSON.parse(localStorage.getItem('apis22')!);
        if (localStorage.getItem('bridges22')) data.bridges = JSON.parse(localStorage.getItem('bridges22')!);
        if (localStorage.getItem('locks22')) data.locks = JSON.parse(localStorage.getItem('locks22')!);

        const apiKeySelect = ref<string>('');
        const bridgeIdSelect = ref<string>('');
        const lockIdSelect = ref<string>('');
        const message = ref<string>('');
        const logs = (): HTMLDivElement => {
            const _div = document.querySelector('.logs');
            if (_div != null) return _div as HTMLDivElement;
            return new HTMLDivElement();
        };
        const redBold = (str: string): string => {
            if (str.indexOf('AT+BATTERY') > -1 || str.indexOf('AT+STATUS') > -1) {
                return `<b style="font-size:18px;color:#f00;">${str}</b>`;
            }
            return str;
        };
        const connWS = () => {
            var host = 'ws://' + location.hostname + ':' + location.port;
            var ws = new WebSocket(`${host}/v2/ttn/websocket?token=${apiKeySelect.value}&applicationId=${bridgeIdSelect.value}&deviceId=${lockIdSelect.value}`);
            ws.onopen = (evt) => {
                logs().innerHTML = `<p class="time">WebSocket Open ${getTime()}</p>` + `<p style="color:red">${JSON.stringify(evt)}</p>` + logs().innerHTML;
            };
            ws.onerror = (evt) => {
                logs().innerHTML = `<p class="time">WebSocket Error ${getTime()}</p>` + `<p style="color:red">${JSON.stringify(evt)}</p>` + logs().innerHTML;
            };
            ws.onclose = (evt) => {
                logs().innerHTML = `<p class="time">WebSocket Close ${getTime()}</p>` + `<p style="color:red">${JSON.stringify(evt)}</p>` + logs().innerHTML;
            };
            ws.onmessage = (evt) => {
                try {
                    const aMsg = JSON.parse(evt.data);
                    var d = new Date(aMsg.result.time),
                        h = d.getHours(),
                        m = d.getMinutes(),
                        s = d.getSeconds();
                    var _time = `${h < 10 ? '0' + h : h}:${m < 10 ? '0' + m : m}:${s < 10 ? '0' + s : s}`;
                    const infoStr = `<p style="color:green">${_time} ${aMsg.result.identifiers[0].device_ids.dev_eui}:${aMsg.result.identifiers[0].device_ids.dev_addr} FPort=${
                        aMsg.result.data.uplink_message.f_port
                    } FrmPayload=${redBold(atob(aMsg.result.data.uplink_message.frm_payload))}</p>`;
                    logs().innerHTML = infoStr + logs().innerHTML;
                } catch {
                    logs().innerHTML = `<p class="time">WebSocket Message ${getTime()}</p>` + `<p style="color:red">${JSON.stringify(evt.data)}</p>` + logs().innerHTML;
                }
            };
            setInterval(() => {
                ws.send('');
            }, 1e3 * 60);
        };
        return {
            saveKey: () => {
                data.apis.push(apiKey.value);
                localStorage.setItem('apis22', JSON.stringify(data.apis));
            },
            saveBridge: () => {
                data.bridges.push(bridgeId.value);
                localStorage.setItem('bridges22', JSON.stringify(data.bridges));
            },
            saveLock: () => {
                data.locks.push(lockId.value);
                localStorage.setItem('locks22', JSON.stringify(data.locks));
            },
            apiKey,
            bridgeId,
            lockId,
            apiKeySelect,
            bridgeIdSelect,
            lockIdSelect,
            data,
            message,
            connWS,
            clearLogs: () => {
                logs().innerText = '';
            },
            sendMessage: () => {
                const t1 = getTime();
                axios
                    .post('/v2/ttn/sendMessage', {
                        apiKey: apiKeySelect.value,
                        applicationId: bridgeIdSelect.value,
                        deviceId: lockIdSelect.value,
                        data: btoa(message.value)
                    })
                    .then((data) => {
                        logs().innerHTML = `<p class="time">sendMessage ${t1}-${getTime()}</p>` + data.data + logs().innerHTML;
                    })
                    .catch((err) => {
                        logs().innerHTML = `<p class="time">sendMessage ${t1}-${getTime()}</p>` + `<p style="color:red">${JSON.stringify(err)}</p>` + logs().innerHTML;
                    });
            },
            fnLock: () => {
                const t1 = getTime();
                axios
                    .post('/v2/ttn/sendMessage', {
                        apiKey: apiKeySelect.value,
                        applicationId: bridgeIdSelect.value,
                        deviceId: lockIdSelect.value,
                        data: btoa('+LOCK\r\n')
                    })
                    .then((data) => {
                        logs().innerHTML = `<p class="time">LOCK ${t1}-${getTime()}</p>` + data.data + logs().innerHTML;
                    })
                    .catch((err) => {
                        logs().innerHTML = `<p class="time">LOCK ${t1}-${getTime()}</p>` + `<p style="color:red">${JSON.stringify(err)}</p>` + logs().innerHTML;
                    });
            },
            fnUnlock: () => {
                const t1 = getTime();
                axios
                    .post('/v2/ttn/sendMessage', {
                        apiKey: apiKeySelect.value,
                        applicationId: bridgeIdSelect.value,
                        deviceId: lockIdSelect.value,
                        data: btoa('+UNLOCK\r\n')
                    })
                    .then((data) => {
                        logs().innerHTML = `<p class="time">UNLOCK ${t1}-${getTime()}</p>` + data.data + logs().innerHTML;
                    })
                    .catch((err) => {
                        logs().innerHTML = `<p class="time">UNLOCK ${t1}-${getTime()}</p>` + `<p style="color:red">${JSON.stringify(err)}</p>` + logs().innerHTML;
                    });
            },
            fnBattery: () => {
                const t1 = getTime();
                axios
                    .post('/v2/ttn/sendMessage', {
                        apiKey: apiKeySelect.value,
                        applicationId: bridgeIdSelect.value,
                        deviceId: lockIdSelect.value,
                        cmd: 'battery',
                        data: btoa('+BATTERY\r\n')
                    })
                    .then((data) => {
                        logs().innerHTML = `<p class="time">BATTERY ${t1}-${getTime()}</p>` + data.data + logs().innerHTML;
                    })
                    .catch((err) => {
                        logs().innerHTML = `<p class="time">BATTERY ${t1}-${getTime()}</p>` + `<p style="color:red">${JSON.stringify(err)}</p>` + logs().innerHTML;
                    });
            },
            fnState: () => {
                const t1 = getTime();
                axios
                    .post('/v2/ttn/sendMessage', {
                        apiKey: apiKeySelect.value,
                        applicationId: bridgeIdSelect.value,
                        deviceId: lockIdSelect.value,
                        cmd: 'state',
                        data: btoa('+STATUS\r\n')
                    })
                    .then((data) => {
                        logs().innerHTML = `<p class="time">STATUS ${t1}-${getTime()}</p>` + data.data + logs().innerHTML;
                    })
                    .catch((err) => {
                        logs().innerHTML = `<p class="time">STATUS ${t1}-${getTime()}</p>` + `<p style="color:red">${JSON.stringify(err)}</p>` + logs().innerHTML;
                    });
            },
            fnMessage: () => {
                const t1 = getTime();
                axios
                    .post('/v2/ttn/message', {
                        apiKey: apiKeySelect.value,
                        applicationId: bridgeIdSelect.value,
                        deviceId: lockIdSelect.value,
                        data: btoa('AT+STATUS')
                    })
                    .then((data) => {
                        logs().innerHTML = `<p class="time">STATUS ${t1}-${getTime()}</p>` + data.data + logs().innerHTML;
                    })
                    .catch((err) => {
                        logs().innerHTML = `<p class="time">STATUS ${t1}-${getTime()}</p>` + `<p style="color:red">${JSON.stringify(err)}</p>` + logs().innerHTML;
                    });
            }
        };
    }
});
</script>

<style lang="less" scoped>
div.main {
    width: 100%;
    max-width: 800px;
    margin: 0 auto;
}
div.model {
    padding: 5px;
    margin: 10px;
    overflow: hidden;
    .m-left {
        width: 60%;
        float: left;
    }
    .m-right {
        width: 35%;
        float: right;
    }
    .v-radio {
        margin: 0 15px;
    }
    .v-btn {
        margin: 5px 1.66667%;
        width: 30%;
    }
}
</style>
<style lang="less">
div.logs {
    padding: 5px;
    margin: 10px 0;
    overflow: hidden;
    width: 100%;
    border: solid 2px #ccc;
    word-break: break-all;
    > div,
    > p {
        background-color: #fefefe;
        margin: 2px 0;
        padding: 5px;
        font-size: 13px;
        line-height: 1.1em;
    }
    > p.time {
        color: #000;
        font-size: 12px;
        padding: 0 5px;
        margin: 0;
    }
}
</style>