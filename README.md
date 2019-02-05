# apicurio-operator

Deploy and manages apicurio by processing its openshift template in a kubernetes/openshift operator.

## Usage

Apply the operator service accounts and its permissions:

```
oc apply -f deploy/role.yaml -f deploy/role_binding.yaml -f deploy/service_account.yaml
```

Create the required crds (needs cluster admin permission):

```
oc apply -f deploy/crds/integreatly_v1alpha1_apicuriodeployment_crd.yaml
```

Deploy the operator:

```
oc apply -f deploy/operator.yaml
```

Create a cr (pick one of your choice):

### Standalone Instance

Deploys a full apicurio instance.

Required parameters:
* openshift app subdomain host

```
oc apply -f deploy/crds/integreatly_v1alpha1_apicuriodeployment_cr.full.yaml
```

### External Keycloak

Deploys apicurio but skips `apicurio-studio-auth` deployment in favor to use an external keycloak instance.

Required parameters:
* openshift app subdomain host
* keycloak host
* keycloak username
* keylcloak password

NOTE: Currently the external keycloak instance requires an pre-configured apicurio realm, instructions can be found here: https://apicurio-studio.readme.io/docs/setting-up-keycloak-for-use-with-apicurio

```
oc apply -f deploy/crds/integreatly_v1alpha1_apicuriodeployment_cr.external_kc.yaml
```

You can check your newly created CR by acessing openshift console in `Resources -> Custom Resources -> Select Apicurio` or by running the command:

```
oc get apicurio
```

Chaning an apicurio cr property, such as the api server memory limit, will trigger a new deployment of the api server.

Deleting a CR will delete all apicurio related resources as well as the apicurio CR is supposed to be a represenation of an apicurio instance.