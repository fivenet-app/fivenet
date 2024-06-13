<!-- eslint-disable vue/multi-word-component-names -->
<script lang="ts" setup>
import { useSettingsStore } from '~/store/settings';

defineOptions({
    inheritAttrs: false,
});

const settingsStore = useSettingsStore();
const { audio: audioSettings } = storeToRefs(settingsStore);

const soundSys = useSound();

const sounds = useState<Sound[]>('sounds', () => []);
</script>

<template>
    <Teleport to="body">
        <div v-if="sounds.length">
            <div v-for="sound of sounds" :key="sound.id">
                <LazyOverlaysSound
                    :id="sound.id"
                    :name="sound.name"
                    :volume="sound.volume ?? audioSettings.notificationsVolume"
                    :rate="sound.rate"
                    @close="soundSys.stop(sound.id!)"
                />
            </div>
        </div>
    </Teleport>
</template>
