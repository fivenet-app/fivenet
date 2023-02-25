<script lang="ts">
import { defineComponent } from 'vue'
import * as grpcWeb from 'grpc-web';
import * as EchoPB from "@arpanet/gen/echo/echo_pb";
import { EchoServiceClient } from "@arpanet/gen/echo/EchoServiceClientPb";

const echoService = new EchoServiceClient('http://localhost:8080', null, null);

export default defineComponent({
    name: "app",
    components: {},
    data: function () {
        return {
            inputField: "" as string,
            message: "" as string,
        };
    },
    created: function () {
        this.echo();
    },
    methods: {
        echo: function () {
            const request = new EchoPB.EchoRequest();
            request.setMessage('Hello World ' + Math.random() + '!');
            echoService.echo(request, { 'custom-header-1': 'value1' },
                (err: grpcWeb.RpcError, response: EchoPB.EchoResponse) => {
                    alert(response.getMessage());
                });
        },
    }
});

</script>

<template>
    <div class="card">
        <button type="button" @click=echo()>message is {{ message }}</button>
        <p>
            Edit
            <code>components/HelloWorld.vue</code> to test HMR
        </p>
    </div>
</template>
