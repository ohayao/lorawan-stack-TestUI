<template>
    <div class="main">
        <div class="model">
            <v-input class="m-left" type="success" v-model:value="apiKey"></v-input>
            <v-btn class="m-right" type="success" size="large" @click="saveKey">添加 API Key</v-btn>
        </div>
        <div class="model">
            <v-input class="m-left" type="pink" v-model:value="bridgeId"></v-input>
            <v-btn class="m-right" type="pink" size="large" @click="saveBridge">添加 Bridge</v-btn>
        </div>
        <div class="model">
            <v-input class="m-left" type="purple" v-model:value="lockId"></v-input>
            <v-btn class="m-right" type="purple" size="large" @click="saveLock">添加 Lock</v-btn>
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
            <v-btn type="primary" size="large" @click="postLock">Lock</v-btn>
            <v-btn type="warn" size="large" @click="postUnlock">Unlock</v-btn>
            <v-btn type="success" size="large" @click="postGetStatus">GetStatus</v-btn>
            <v-btn type="pink" size="large" @click="getBridgeLogs">BridgeLogs</v-btn>
            <v-btn type="purple" size="large" @click="addLock">绑定锁</v-btn>
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
        if (localStorage.getItem('apis')) data.apis = JSON.parse(localStorage.getItem('apis')!);
        if (localStorage.getItem('bridges')) data.bridges = JSON.parse(localStorage.getItem('bridges')!);
        if (localStorage.getItem('locks')) data.locks = JSON.parse(localStorage.getItem('locks')!);

        const apiKeySelect = ref<string>('');
        const bridgeIdSelect = ref<string>('');
        const lockIdSelect = ref<string>('');
        console.log(data);
        const logs = (): HTMLDivElement => {
            const _div = document.querySelector('.logs');
            if (_div != null) return _div as HTMLDivElement;
            return new HTMLDivElement();
        };
        return {
            saveKey: () => {
                data.apis.push(apiKey.value);
                localStorage.setItem('apis', JSON.stringify(data.apis));
            },
            saveBridge: () => {
                data.bridges.push(bridgeId.value);
                localStorage.setItem('bridges', JSON.stringify(data.bridges));
            },
            saveLock: () => {
                data.locks.push(lockId.value);
                localStorage.setItem('locks', JSON.stringify(data.locks));
            },
            apiKey,
            bridgeId,
            lockId,
            apiKeySelect,
            bridgeIdSelect,
            lockIdSelect,
            data,
            clearLogs: () => {
                logs().innerText = '';
            },
            postLock: () => {
                const t1 = getTime();
                axios
                    .post('/v1/api/lock', {
                        apiKey: apiKeySelect.value,
                        bridgeId: bridgeIdSelect.value,
                        lockId: lockIdSelect.value,
                        expiryDuration: '5'
                    })
                    .then((data) => {
                        logs().innerHTML = `<p class="time">Lock ${t1}-${getTime()}</p>` + data.data + logs().innerHTML;
                    })
                    .catch((err) => {
                        logs().innerHTML = `<p class="time">Lock ${t1}-${getTime()}</p>` + `<p style="color:red">${JSON.stringify(err)}</p>` + logs().innerHTML;
                    });
            },
            postUnlock: () => {
                const t1 = getTime();
                axios
                    .post('/v1/api/unlock', {
                        apiKey: apiKeySelect.value,
                        bridgeId: bridgeIdSelect.value,
                        lockId: lockIdSelect.value,
                        expiryDuration: '5'
                    })
                    .then((data) => {
                        logs().innerHTML = `<p class="time">Unlock ${t1}-${getTime()}</p>` + data.data + logs().innerHTML;
                    })
                    .catch((err) => {
                        logs().innerHTML = `<p class="time">Unlock ${t1}-${getTime()}</p>` + `<p style="color:red">${JSON.stringify(err)}</p>` + logs().innerHTML;
                    });
            },
            postGetStatus: () => {
                const t1 = getTime();
                axios
                    .post('/v1/api/getstatus', {
                        apiKey: apiKeySelect.value,
                        bridgeId: bridgeIdSelect.value,
                        lockId: lockIdSelect.value,
                        expiryDuration: '5'
                    })
                    .then((data) => {
                        logs().innerHTML = `<p class="time">GetStatus ${t1}-${getTime()}</p>` + data.data + logs().innerHTML;
                    })
                    .catch((err) => {
                        logs().innerHTML = `<p class="time">GetStatus ${t1}-${getTime()}</p>` + `<p style="color:red">${JSON.stringify(err)}</p>` + logs().innerHTML;
                    });
            },
            getBridgeLogs: () => {
                const t1 = getTime();
                axios
                    .post('/v1/api/getbridgelogs', {
                        apiKey: apiKeySelect.value,
                        bridgeId: bridgeIdSelect.value
                    })
                    .then((data) => {
                        logs().innerHTML = `<p class="time">GetLogs ${t1}-${getTime()}</p>` + data.data + logs().innerHTML;
                    })
                    .catch((err) => {
                        logs().innerHTML = `<p class="time">GetLogs ${t1}-${getTime()}</p>` + `<p style="color:red">${JSON.stringify(err)}</p>` + logs().innerHTML;
                    });
            },
            addLock: () => {
                const t1 = getTime();
                axios
                    .post('/v1/api/addlock', {
                        apiKey: apiKeySelect.value,
                        bridgeId: bridgeIdSelect.value,
                        lockId: lockIdSelect.value
                    })
                    .then((data) => {
                        logs().innerHTML = `<p class="time">AddLock ${t1}-${getTime()}</p>` + data.data + logs().innerHTML;
                    })
                    .catch((err) => {
                        logs().innerHTML = `<p class="time">AddLock ${t1}-${getTime()}</p>` + `<p style="color:red">${JSON.stringify(err)}</p>` + logs().innerHTML;
                    });
            }
        };
    }
});
</script>

<style lang="less" scoped>
div.main {
    width: 100%;
    max-width: 640px;
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
    > div,
    > p {
        background-color: #fefefe;
        margin: 2px 0;
        padding: 5px;
        font-size: 13px;
        line-height: 1.1em;
        &:hover {
            background-color: #000;
            color: #fff;
        }
    }
    > p.time {
        color: #000;
        font-size: 12px;
        padding: 0 5px;
        margin: 0;
        &:hover {
            background-color: #fefefe;
            color: #000;
        }
    }
}
</style>