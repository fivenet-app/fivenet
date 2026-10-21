import type { Ref } from 'vue';
import {
    PropertyDefinition,
    PropertyDefinitionSource,
    PropertyDisplayMode,
    PropertyRenderer,
    PropertyScope,
    PropertyValueType,
    type Property,
} from '~~/gen/ts/resources/properties/properties';
import type { User } from '~~/gen/ts/resources/users/user';

export function useCitizenProperties(user: Ref<User>) {
    const { attr } = useAuth();

    const canViewPhone = attr('citizens.CitizensService/ListCitizens', 'Properties', 'user.phone_number');
    const canViewVisum = attr('citizens.CitizensService/ListCitizens', 'Properties', 'user.visum');
    const canViewEmail = attr('citizens.CitizensService/ListCitizens', 'Properties', 'user.email');

    const coreProperties = computed<Property[]>(() => {
        const result: Property[] = [];
        const definition = (
            key: string,
            displayKey: string,
            valueType: PropertyValueType,
            groupKey = 'personal',
            renderer = PropertyRenderer.DEFAULT,
        ) =>
            PropertyDefinition.create({
                scope: PropertyScope.USER,
                key,
                valueType,
                source: PropertyDefinitionSource.SYSTEM,
                enabled: true,
                displayKey,
                groupKey,
                sortOrder: 0,
                renderer,
                displayMode: PropertyDisplayMode.DEFAULT,
            });

        if (user.value.dateofbirth) {
            result.push({
                key: 'user.dateofbirth',
                value: { value: { oneofKind: 'stringValue', stringValue: user.value.dateofbirth } },
                definition: definition('user.dateofbirth', 'common.date_of_birth', PropertyValueType.STRING),
            });
        }
        if (user.value.sex) {
            result.push({
                key: 'user.sex',
                value: { value: { oneofKind: 'stringValue', stringValue: user.value.sex } },
                definition: definition('user.sex', 'common.sex', PropertyValueType.STRING),
            });
        }
        if (user.value.height !== undefined) {
            result.push({
                key: 'user.height',
                value: { value: { oneofKind: 'numberValue', numberValue: user.value.height } },
                definition: definition(
                    'user.height',
                    'common.height',
                    PropertyValueType.NUMBER,
                    'personal',
                    PropertyRenderer.HEIGHT,
                ),
            });
        }
        if (canViewPhone.value && user.value.phoneNumber) {
            result.push({
                key: 'user.phone_number',
                value: { value: { oneofKind: 'stringValue', stringValue: user.value.phoneNumber } },
                definition: definition(
                    'user.phone_number',
                    'common.phone_number',
                    PropertyValueType.STRING,
                    'contact',
                    PropertyRenderer.PHONE,
                ),
            });
        }
        if (canViewEmail.value) {
            result.push({
                key: 'user.email',
                value: user.value.props?.email
                    ? { value: { oneofKind: 'stringValue', stringValue: user.value.props.email } }
                    : undefined,
                definition: definition(
                    'user.email',
                    'common.mail',
                    PropertyValueType.STRING,
                    'contact',
                    PropertyRenderer.EMAIL,
                ),
            });
        }
        if (canViewVisum.value && user.value.visum !== undefined) {
            result.push({
                key: 'user.visum',
                value: { value: { oneofKind: 'integerValue', integerValue: user.value.visum } },
                definition: definition('user.visum', 'common.visum', PropertyValueType.INTEGER, 'contact'),
            });
        }
        return result;
    });

    const profileProperties = computed<Property[]>(() => {
        const dynamicProperties = user.value.properties?.properties ?? [];
        const dynamicKeys = new Set(dynamicProperties.map((property) => property.key));
        return [...coreProperties.value.filter((property) => !dynamicKeys.has(property.key)), ...dynamicProperties];
    });

    return { coreProperties, profileProperties };
}
