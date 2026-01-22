/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"testing"
)

// TestValidateDescription makes sure the validation regexp  allows the
// expected ranges of characters.
//
// Most of the test cases have been built using the descriptions taken from
// images-private/images/*/manifest.yaml Here is the one-liner that got the
// data:
//
// fd -HI -t f '^metadata\.yaml$' | xargs -I{} yq -r '. as $m | "{\nName: \"\($m.name)\",\nInput: \"\($m.short_description)\",\nExpect: true,\n},"' "{}"
func TestValidateDescription(t *testing.T) {
	tests := []struct {
		Name   string
		Input  string
		Expect bool
	}{
		{
			Name:   "empty",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "newline isn't allowed",
			Input:  "this is\ndescription",
			Expect: false,
		},
		{
			Name:   "tab isn't allowed",
			Input:  "this is\tdescription",
			Expect: false,
		},
		{
			Name:   "carriage-return isn't allowed",
			Input:  "this is\rdescription",
			Expect: false,
		},
		{
			Name:   "harbor-exporter-fips",
			Input:  "A Wolf-based image for Harbor Exporter - application for monitoring harbor deployments.",
			Expect: true,
		},
		{
			Name:   "postgres-iamguarded",
			Input:  "PostgreSQL is a powerful, open source object-relational database system.",
			Expect: true,
		},
		{
			Name:   "request-6212",
			Input:  "PostgreSQL IAMGuarded FIPS with pgvector",
			Expect: true,
		},
		{
			Name:   "flux-helm-controller-iamguarded-fips",
			Input:  "minimal zero CVE flux images",
			Expect: true,
		},
		{
			Name:   "influxdb-iamguarded",
			Input:  "[InfluxDB](https://github.com/influxdata/influxdb) is a time series database designed to handle high write and query loads.",
			Expect: true,
		},
		{
			Name:   "strongswan",
			Input:  "strongSwan is an OpenSource IPsec-based VPN solution",
			Expect: true,
		},
		{
			Name:   "tritonserver-fips",
			Input:  "The Triton Inference Server provides an optimized cloud and edge inferencing solution.",
			Expect: true,
		},
		{
			Name:   "buildkit",
			Input:  "Buildkit is a concurrent, cache-efficient, and Dockerfile-agnostic builder toolkit.",
			Expect: true,
		},
		{
			Name:   "jre",
			Input:  "Minimalist Wolfi-based Java JRE image using [OpenJDK](https://openjdk.org/projects/jdk/). Used for running Java applications.",
			Expect: true,
		},
		{
			Name:   "php-fpm_exporter",
			Input:  "Minimal php-fpm_exporter Image",
			Expect: true,
		},
		{
			Name:   "vllm-openai",
			Input:  "vLLM is a high-throughput and memory-efficient inference engine for Large Language Models (LLMs). It provides an OpenAI-compatible API server for production LLM deployments with GPU acceleration.",
			Expect: true,
		},
		{
			Name:   "ant",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "vt",
			Input:  "Minimal image with the Virus Total CLI - `vt-cli`.",
			Expect: true,
		},
		{
			Name:   "yunikorn-scheduler",
			Input:  "Apache YuniKorn K8shim",
			Expect: true,
		},
		{
			Name:   "aws-load-balancer-controller",
			Input:  "Minimal Image for Kubernetes controller for Elastic Load Balancers",
			Expect: true,
		},
		{
			Name:   "step-ca-fips",
			Input:  "Minimal FIPS image of [step-ca](https://smallstep.com/docs/step-ca), an online Certificate Authority (CA) for secure, automated X.509 and SSH certificate management",
			Expect: true,
		},
		{
			Name:   "rabbitmq-fips",
			Input:  "[RabbitMQ](https://github.com/rabbitmq/rabbitmq-server) is a message broker.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "neo4j",
			Input:  "Neo4J is a graph database that is commonly used in applications that require complex relationships between data. Neo4j supports both a standalone and a cluster deployment of Neo4j on Kubernetes using the Neo4j Helm charts.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "promtail-fips",
			Input:  "This image contains the `promtail` application for log aggregation. `promtail` is the log aggregator that ships logs to Loki and/or Prometheus. It runs as an agent and scrapes logs from files, containers, and hosts and ships them to a logging backend.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "prometheus-statsd-exporter-fips",
			Input:  "Minimalist Wolfi-based Prometheus StatsD Exporter image for exporting metrics to StatsD.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "dragonfly-operator-fips",
			Input:  "A FIPS-compliant Kubernetes operator used to deploy and manage Dragonfly instances inside your Kubernetes clusters.",
			Expect: true,
		},
		{
			Name:   "spamcheck-fips",
			Input:  "Spamcheck is a gRPC-based spam classification service for GitLab",
			Expect: true,
		},
		{
			Name:   "kube-state-metrics-iamguarded",
			Input:  "Kube-state-metrics generates Prometheus metrics about Kubernetes objects",
			Expect: true,
		},
		{
			Name:   "syft-fips",
			Input:  "A tool for generating a Software Bill of Materials (SBOM) from container images and filesystems.",
			Expect: true,
		},
		{
			Name:   "grafana-alloy-operator",
			Input:  "The Alloy Operator is a Kubernetes Operator that manages the lifecycle of Grafana Alloy instances",
			Expect: true,
		},
		{
			Name:   "selenium",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "spire-controller-manager",
			Input:  "The SPIRE Controller Manager provides automated workload identity management for Kubernetes clusters through SPIRE",
			Expect: true,
		},
		{
			Name:   "dogstatsd",
			Input:  "Standalone DogStatsD image for custom metrics collection",
			Expect: true,
		},
		{
			Name:   "fluent-bit-plugin-loki",
			Input:  "The Fluent Bit Loki plugin allows you to send your log or events to a Loki service.",
			Expect: true,
		},
		{
			Name:   "vela-cli",
			Input:  "Reference implementation for Kubevela's vela CLI tool.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "rabbitmq-messaging-topology-operator-iamguarded",
			Input:  "RabbitMQ messaging topology operator",
			Expect: true,
		},
		{
			Name:   "haproxy-iamguarded-fips",
			Input:  "A minimal [haproxy](https://www.haproxy.org/) base image rebuilt every night from source.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "frr",
			Input:  "The FRRouting Protocol Suite",
			Expect: true,
		},
		{
			Name:   "cert-manager-csi-driver-fips",
			Input:  "A Kubernetes CSI driver that automatically mounts signed certificates to Pods using ephemeral volumes",
			Expect: true,
		},
		{
			Name:   "terraform",
			Input:  "[Terraform](https://github.com/hashicorp/terraform) is an infrastructure as code tool.",
			Expect: true,
		},
		{
			Name:   "mountpoint-s3-csi-driver",
			Input:  "Built on Mountpoint for Amazon S3, the Mountpoint CSI driver presents an Amazon S3 bucket as a storage volume accessible by containers in your Kubernetes cluster.",
			Expect: true,
		},
		{
			Name:   "multus-cni",
			Input:  "A CNI meta-plugin for multi-homed pods in Kubernetes",
			Expect: true,
		},
		{
			Name:   "prometheus-logstash-exporter",
			Input:  "Prometheus exporter for Logstash written in Go",
			Expect: true,
		},
		{
			Name:   "librechat-fips",
			Input:  "A FIPS-compliant AI chat application.",
			Expect: true,
		},
		{
			Name:   "mongodb-iamguarded",
			Input:  "[MongoDB](https://www.mongodb.com/) is a document-oriented database management system. MongoDB is a popular example of a NoSQL database, and stores data in JSON-like documents.",
			Expect: true,
		},
		{
			Name:   "memcached-exporter-iamguarded-fips",
			Input:  "A memcached exporter for Prometheus.",
			Expect: true,
		},
		{
			Name:   "adoptium-jdk",
			Input:  "Minimal Wolfi-based Adoptium Java JDK image using [Adoptium OpenJDK](https://adoptium.net).  Used for compiling Java applications.",
			Expect: true,
		},
		{
			Name:   "external-secrets-fips",
			Input:  "Fetches secrets from external systems and exposes them as Kubernetes Secrets.",
			Expect: true,
		},
		{
			Name:   "prometheus-blackbox-exporter-iamguarded",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "keycloak-config-cli-iamguarded-fips",
			Input:  "Import YAML/JSON-formatted configuration files into Keycloak - Configuration as Code for Keycloak.",
			Expect: true,
		},
		{
			Name:   "openresty-fips",
			Input:  "OpenResty is a high Performance Web Platform Based on Nginx and LuaJIT.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "kubernetes-dashboard-metrics-scraper-fips",
			Input:  "Module containing the Kubernetes metrics scraper module of the Kubernetes dashboard application",
			Expect: true,
		},
		{
			Name:   "mailpit",
			Input:  "Mailpit is an email and SMTP testing tool with API for developers.",
			Expect: true,
		},
		{
			Name:   "yunikorn-web",
			Input:  "Apache YuniKorn Web UI",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "aspnet-runtime",
			Input:  "Container image with the latest ASP.NET runtime.",
			Expect: true,
		},
		{
			Name:   "ollama-fips",
			Input:  "Get up and running with Llama 3.3, DeepSeek-R1, Phi-4, Gemma 3, Mistral Small 3.1 and other large language models.",
			Expect: true,
		},
		{
			Name:   "rabbitmq-default-user-credential-updater-iamguarded-fips",
			Input:  "Image with [default-user-credential-updater](https://github.com/rabbitmq/default-user-credential-updater)",
			Expect: true,
		},
		{
			Name:   "dbgate-fips",
			Input:  "FIPS-hardened version of DbGate, a database administration tool for SQL Server, MySQL, PostgreSQL, MongoDB, Redis and SQLite with FIPS 140-2 compliant cryptographic operations.",
			Expect: true,
		},
		{
			Name:   "grafana-beyla",
			Input:  "Open source eBPF-based auto-instrumentation tool that helps you easily get started with application observability",
			Expect: true,
		},
		{
			Name:   "pushprox-fips",
			Input:  "Minimal FIPS compliant image with [PushProx](https://github.com/prometheus-community/PushProx), a proxy to allow Prometheus to scrape through NAT etc.",
			Expect: true,
		},
		{
			Name:   "pdns-recursor",
			Input:  "PowerDNS Recursor is a non authoritative/recursing DNS server.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "tritonserver-vllm-backend",
			Input:  "The Triton Inference Server provides an optimized cloud and edge inferencing solution with vllm backend",
			Expect: true,
		},
		{
			Name:   "nginx-prometheus-exporter-fips",
			Input:  "The `nginx-prometheus-exporter-fips` image is designed to scrape metrics from an NGINX instance and expose them to Prometheus in a secure and minimal environment. Below are detailed instructions for using the image in both Docker and Kubernetes environments.",
			Expect: true,
		},
		{
			Name:   "loki",
			Input:  "This image contains the `loki` application for log aggregation. `loki` can be used to stream, aggregate, and query logs from apps and infrastructure.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "Minimalist Wolfi-based Java JRE image using Adoptium OpenJDK. Used for running Java applications.",
			Expect: true,
		},
		{
			Name:   "gatekeeper-crds",
			Input:  "Minimal image for installing Gatekeeper Custom Resource Definitions (CRDs) in Kubernetes clusters.",
			Expect: true,
		},
		{
			Name:   "prometheus-mongodb-exporter-iamguarded",
			Input:  "Prometheus MongoDB Exporter image for exporting various metrics about MongoDB compatible with Iamguarded Helm chart.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "crossplane-provider-terraform-fips",
			Input:  "Minimal Fips image of Crossplane Terraform Provider enables provisioning and managing infrastructure using Terraform within a Crossplane control plane.",
			Expect: true,
		},
		{
			Name:   "kubernetes-csi-external-health-monitor",
			Input:  "CSI External Health Monitor Controller",
			Expect: true,
		},
		{
			Name:   "opentelemetry-collector",
			Input:  "Minimal image with [opentelemetry-collector](https://github.com/open-telemetry/opentelemetry-collector)",
			Expect: true,
		},
		{
			Name:   "k8s-wait-for-fips",
			Input:  "Container image for waiting for a k8s service, job or pods to enter a desired state.",
			Expect: true,
		},
		{
			Name:   "prometheus-redis-exporter-iamguarded-fips",
			Input:  "Prometheus Redis Exporter image for exporting metrics to Redis.",
			Expect: true,
		},
		{
			Name:   "chrony_exporter",
			Input:  "Minimalist Wolfi-based image for Chrony Exporter - a Prometheus exporter for Chrony NTP metrics.",
			Expect: true,
		},
		{
			Name:   "prometheus-blackbox-exporter-iamguarded-fips",
			Input:  "Prometheus blackbox exporter allows blackbox probing of endpoints over HTTP, HTTPS, DNS, TCP, ICMP and gRPC.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "images/emqx-exporter-fips",
			Input:  "The emqx-exporter is designed to expose partial metrics that are not included in the EMQX Prometheus API.",
			Expect: true,
		},
		{
			Name:   "thanos-fips",
			Input:  "Minimal Thanos Image, a highly available Prometheus setup with long term storage",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "dnsdist-fips",
			Input:  "dnsdist is a highly DNS-, DoS- and abuse-aware loadbalancer",
			Expect: true,
		},
		{
			Name:   "prometheus-pushgateway-iamguarded",
			Input:  "Push acceptor for ephemeral and batch jobs.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "minio-iamguarded-fips",
			Input:  "MinIO is a high-performance, S3 compatible object store. This FIPS-compliant iamguarded variant is specifically designed to work with the iamguarded Helm chart in FIPS-140 compliant environments.",
			Expect: true,
		},
		{
			Name:   "opensearch-k8s-operator-fips",
			Input:  "The FIPS-complaint OpenSearch Kubernetes Operator is used for automating the deployment, provisioning, management, and orchestration of OpenSearch clusters and OpenSearch dashboards.",
			Expect: true,
		},
		{
			Name:   "openbao",
			Input:  "Minimal image with OpenBao.",
			Expect: true,
		},
		{
			Name:   "flannel-cni-plugin-fips",
			Input:  "A plugin designed to work in conjunction with flannel, a network fabric for containers.",
			Expect: true,
		},
		{
			Name:   "sbt-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "sealed-secrets-kubeseal-fips",
			Input:  "A Kubernetes tool for one-way encrypted Secrets, enabling safe GitOps-friendly secret management.",
			Expect: true,
		},
		{
			Name:   "wazero",
			Input:  "This image contains the `wazero` tool which can be used to compile or run wasm binaries.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "nginx-prometheus-exporter-iamguarded",
			Input:  "Nginx prometheus exporter",
			Expect: true,
		},
		{
			Name:   "prometheus-postgres-exporter-iamguarded-fips",
			Input:  "A PostgreSQL metric exporter for Prometheus",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "nextflow-fips",
			Input:  "Nextflow is a domain-specific language (DSL) for data-driven computational pipelines.",
			Expect: true,
		},
		{
			Name:   "cloud-provider-aws-fips",
			Input:  "Cloud provider for AWS",
			Expect: true,
		},
		{
			Name:   "cloudnative-pg",
			Input:  "Minimal Wolfi-based image for CloudNative PG, a comprehensive platform designed to seamlessly manage PostgreSQL databases within Kubernetes environments.",
			Expect: true,
		},
		{
			Name:   "etcd-iamguarded",
			Input:  "[etcd](https://etcd.io/) Distributed reliable key-value store for the most critical data of a distributed system",
			Expect: true,
		},
		{
			Name:   "neuvector-manager",
			Input:  "A Wolfi-based image for NeuVector - a full lifecycle container security platform.",
			Expect: true,
		},
		{
			Name:   "clickhouse-keeper-fips",
			Input:  "ClickHouse Keeper is a distributed coordination service that provides a ZooKeeper-compatible API for managing ClickHouse clusters. It handles distributed consensus, configuration management, and leader election using the Raft algorithm.",
			Expect: true,
		},
		{
			Name:   "k6-fips",
			Input:  "Load testing tool for testing APIs, microservices, and websites.",
			Expect: true,
		},
		{
			Name:   "kiali",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "stakater-reloader-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "cc-dynamic",
			Input:  "Base image with just enough to run arbitrary binaries that may require gcc or cc libraries, typically C++ or Rust binaries.",
			Expect: true,
		},
		{
			Name:   "secrets-store-csi-driver-fips",
			Input:  "Minimal image with Kubernetes Secrets Store CSI Driver.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "k8s-metadata-injection-fips",
			Input:  "Kubernetes metadata injection for New Relic APM to make a linkage between APM and Infrastructure data. This is the FIPS-compliant variant of the image.",
			Expect: true,
		},
		{
			Name:   "ipfs-cluster",
			Input:  "Pinset orchestration for IPFS",
			Expect: true,
		},
		{
			Name:   "promxy",
			Input:  "Minimal image with Promxy.",
			Expect: true,
		},
		{
			Name:   "knative-net-istio-webhook",
			Input:  "Istio uses webhook for validating Istio configuration (ValidatingAdmissionWebhooks) and automatically injecting the sidecar proxy into user pods (MutatingAdmissionWebhooks).",
			Expect: true,
		},
		{
			Name:   "doppler-kubernetes-operator",
			Input:  "Automatically sync secrets from Doppler to Kubernetes and auto-reload deployments when secrets change.",
			Expect: true,
		},
		{
			Name:   "kubernetes-fips",
			Input:  "Production-Grade Container Scheduling and Management",
			Expect: true,
		},
		{
			Name:   "dcgm-fips",
			Input:  "NVIDIA Data Center GPU Manager (DCGM) is a project for gathering telemetry and measuring the health of NVIDIA GPUs",
			Expect: true,
		},
		{
			Name:   "timoni",
			Input:  "Minimal image with `timoni` binary. `timoni` is a package manager for Kubernetes, powered by `cue` and inspired by `helm`.",
			Expect: true,
		},
		{
			Name:   "vela-core-fips",
			Input:  "KubeVela is a modern software delivery platform that makes deploying and operating applications across today's hybrid, multi-cloud environments easier, faster and more reliable.",
			Expect: true,
		},
		{
			Name:   "azure-service-operator",
			Input:  "Instead of deploying and managing your Azure resources separately from your Kubernetes application, ASO allows you to manage them together, automatically configuring your application as needed. For example, ASO can set up your Redis Cache or PostgreSQL database server and then configure your Kubernetes application to use them.",
			Expect: true,
		},
		{
			Name:   "valkey-bundle",
			Input:  "Valkey bundle with pre-installed modules for extended functionality.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "request-2682",
			Input:  "A container registry offering secure image storage, access control, scanning, and replication, built as Harbor's customized Docker Distribution registry service.",
			Expect: true,
		},
		{
			Name:   "cortex",
			Input:  "Cortex provides horizontally scalable, highly available, multi-tenant, long term storage for Prometheus.",
			Expect: true,
		},
		{
			Name:   "kubernetes-csi-node-driver-registrar",
			Input:  " Sidecar container that registers a CSI driver with the kubelet using the kubelet plugin registration mechanism. ",
			Expect: true,
		},
		{
			Name:   "eck-operator",
			Input:  "Elastic Cloud on Kubernetes",
			Expect: true,
		},
		{
			Name:   "oauth2-proxy-iamguarded-fips",
			Input:  "[OAuth2 Proxy](https://oauth2-proxy.github.io/oauth2-proxy/) is a reverse proxy that provides authentication with Google, Azure, OpenID Connect and many more identity providers.",
			Expect: true,
		},
		{
			Name:   "ffmpeg",
			Input:  "Minimal image that contains ffmpeg",
			Expect: true,
		},
		{
			Name:   "jdk-crac",
			Input:  "Minimalist Wolfi-based OpenJDK image with [CRaC](https://openjdk.org/projects/crac/) support. Used for compiling Java applications.",
			Expect: true,
		},
		{
			Name:   "spark-operator",
			Input:  "A minimal, Wolfi-based image for Spark Operator. Facilitates the deployment and management of Apache Spark applications in Kubernetes environments.",
			Expect: true,
		},
		{
			Name:   "jupyterhub-k8s-image-awaiter-fips",
			Input:  "JupyterHub Kubernetes Image Awaiter - ensures images are pre-pulled before deployment",
			Expect: true,
		},
		{
			Name:   "airflow-iamguarded",
			Input:  "Apache Airflow is a platform to programmatically author, schedule, and monitor workflows as directed acyclic graphs (DAGs).",
			Expect: true,
		},
		{
			Name:   "jmx-exporter-iamguarded-fips",
			Input:  "jmx-exporter-iamguarded-fips image is a Prometheus metrics exporter for JMX (Java Management Extensions), enabling monitoring and observability for Java applications used by IAMguarded helm charts.",
			Expect: true,
		},
		{
			Name:   "victoriametrics",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "geoip-api",
			Input:  "A JSON REST API for Maxmind GeoIP databases",
			Expect: true,
		},
		{
			Name:   "cloudprober",
			Input:  "Cloudprober is a monitoring software that makes it super-easy to monitor availability and performance of various components of your system",
			Expect: true,
		},
		{
			Name:   "docker-dind",
			Input:  "Chainguard image for Docker in Docker (DinD), allowing you to run Docker within a container.",
			Expect: true,
		},
		{
			Name:   "bats",
			Input:  "Bats provides a simple way to verify that the UNIX programs you write behave as expected.",
			Expect: true,
		},
		{
			Name:   "dragonfly",
			Input:  "Dragonfly is a modern drop-in replacement for Redis and Memcached, offering better performance, multithreading, and lower memory overhead.",
			Expect: true,
		},
		{
			Name:   "knative-serving",
			Input:  "Knative Serving builds on Kubernetes to support deploying and serving of applications and functions as serverless containers.",
			Expect: true,
		},
		{
			Name:   "liberica-jdk-fips",
			Input:  "Free and open source Progressive Java Runtime for modern Java deployments",
			Expect: true,
		},
		{
			Name:   "dex-iamguarded-fips",
			Input:  "dex is a federated OpenID Connect provider.",
			Expect: true,
		},
		{
			Name:   "sealed-secrets-controller-iamguarded",
			Input:  "A Kubernetes controller and tool for one-way encrypted Secrets, enabling safe GitOps-friendly secret management.",
			Expect: true,
		},
		{
			Name:   "listmonk",
			Input:  "High performance, self-hosted, newsletter and mailing list manager with a modern dashboard. Single binary app.",
			Expect: true,
		},
		{
			Name:   "oauth2-proxy-fips",
			Input:  "[OAuth2 Proxy](https://oauth2-proxy.github.io/oauth2-proxy/) is a reverse proxy that provides authentication with Google, Azure, OpenID Connect and many more identity providers.",
			Expect: true,
		},
		{
			Name:   "kyverno-policy-reporter-plugin-trivy",
			Input:  "This Plugin for Policy Reporter brings additional Trivy specific information to the Policy Reporter UI",
			Expect: true,
		},
		{
			Name:   "policy-bot",
			Input:  "A GitHub App that enforces approval policies on pull requests",
			Expect: true,
		},
		{
			Name:   "images/reposilite",
			Input:  "Lightweight and easy-to-use repository management software dedicated for the Maven-based artifacts in the JVM ecosystem ",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "clickhouse",
			Input:  "[Clickhouse](https://clickhouse.com) is the fastest and most resource efficient open-source database for real-time apps and analytics.",
			Expect: true,
		},
		{
			Name:   "kubernetes-reflector-fips",
			Input:  "Kubernetes controller for reflecting ConfigMaps, Secrets, and Certificates across namespaces",
			Expect: true,
		},
		{
			Name:   "dcgm",
			Input:  "NVIDIA Data Center GPU Manager (DCGM) is a project for gathering telemetry and measuring the health of NVIDIA GPUs",
			Expect: true,
		},
		{
			Name:   "cluster-api-azure-controller-fips",
			Input:  "Kubernetes Cluster API provider for Microsoft Azure infrastructure management.",
			Expect: true,
		},
		{
			Name:   "rancher-shell",
			Input:  "rancher-shell is the kubectl and helm installer image for Rancher",
			Expect: true,
		},
		{
			Name:   "argo-rollouts-fips",
			Input:  "Argo Rollouts is a Kubernetes controller and set of CRDs which provide advanced deployment capabilities such as blue-green, canary, canary analysis, experimentation, and progressive delivery features to Kubernetes.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "tofu-controller",
			Input:  "A GitOps OpenTofu and Terraform controller for Flux",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "jmx-exporter-iamguarded",
			Input:  "jmx-exporter-iamguarded image is a Prometheus metrics exporter for JMX (Java Management Extensions), enabling monitoring and observability for Java applications used by IAMguarded helm charts.",
			Expect: true,
		},
		{
			Name:   "argocd-image-updater",
			Input:  "Automatic container image update for Argo CD",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "strimzi-kafka",
			Input:  "Strimzi provides a way to run an Apache Kafka® cluster on Kubernetes or OpenShift in various deployment configurations.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "infinispan-operator",
			Input:  "Kubernetes Operator for Infinispan. Automates deployment, scaling, and lifecycle management of Infinispan distributed in-memory data grid clusters.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "request-5510",
			Input:  "A time-series database for high-performance real-time analytics packaged as a Postgres extension",
			Expect: true,
		},
		{
			Name:   "request-5948",
			Input:  "A custom build of dcgm-exporter-fips without setcap SYS_ADMIN",
			Expect: true,
		},
		{
			Name:   "cert-manager-startupapicheck-fips",
			Input:  "Automatically provision and manage TLS certificates in Kubernetes",
			Expect: true,
		},
		{
			Name:   "go",
			Input:  "Container image for building Go applications.",
			Expect: true,
		},
		{
			Name:   "kube-state-metrics",
			Input:  "Minimal Kube State Metrics Image",
			Expect: true,
		},
		{
			Name:   "kyverno-notation-aws",
			Input:  "Kyverno extension service for Notation and the AWS signer",
			Expect: true,
		},
		{
			Name:   "apache-camel-karavan-devmode",
			Input:  "Development container for Apache Camel Karavan low-code integration platform.",
			Expect: true,
		},
		{
			Name:   "az-iamguarded-fips",
			Input:  "Azure CLI (IAM Guarded, FIPS)",
			Expect: true,
		},
		{
			Name:   "deno",
			Input:  "Minimal container image for running [Deno](https://deno.com/) apps",
			Expect: true,
		},
		{
			Name:   "mesosphere-vsphere-csi-syncer-fips",
			Input:  "vSphere storage Container Storage Interface (CSI) plugin",
			Expect: true,
		},
		{
			Name:   "metaflow-metadata-service",
			Input:  "A minimal, wolfi-based image for Metaflow Metadata Service. Metaflow Metadata Service is a backend service for tracking and managing ML workflow metadata.",
			Expect: true,
		},
		{
			Name:   "prometheus-elasticsearch-exporter-iamguarded-fips",
			Input:  "Minimalist Wolfi-based Prometheus Elasticsearch Exporter image for exporting various metrics about Elasticsearch.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "langfuse",
			Input:  "Langfuse is an open-source observability and analytics platform for LLM applications.",
			Expect: true,
		},
		{
			Name:   "azure-service-operator-fips",
			Input:  "Instead of deploying and managing your Azure resources separately from your Kubernetes application, ASO allows you to manage them together, automatically configuring your application as needed. For example, ASO can set up your Redis Cache or PostgreSQL database server and then configure your Kubernetes application to use them.",
			Expect: true,
		},
		{
			Name:   "objectstorage",
			Input:  "Container Object Storage Interface (COSI) controller and sidecar",
			Expect: true,
		},
		{
			Name:   "zig",
			Input:  "Minimal image with zig binary.",
			Expect: true,
		},
		{
			Name:   "nfpm",
			Input:  "nFPM is Not FPM - a simple deb, rpm, apk, ipk, and arch linux packager written in Go",
			Expect: true,
		},
		{
			Name:   "flyway",
			Input:  "Flyway is a database migration tool to evolve your database schema easily and reliably across all your instances.",
			Expect: true,
		},
		{
			Name:   "sonobuoy-systemd-logs",
			Input:  "This is a simple standalone container that gathers log information from systemd, by chrooting into the node's filesystem and running journalctl.",
			Expect: true,
		},
		{
			Name:   "deck-fips",
			Input:  "deck is a command-line interface for managing Kong Gateway configurations declaratively",
			Expect: true,
		},
		{
			Name:   "gitlab-kubectl",
			Input:  "kubectl is the official CLI tool for managing Kubernetes clusters.",
			Expect: true,
		},
		{
			Name:   "chainguard-base-fips",
			Input:  "Minimal, FIPS-validated image useful as a base for building secure images.",
			Expect: true,
		},
		{
			Name:   "descheduler",
			Input:  "Kubernetes Descheduler is a tool that evicts pods from nodes based on configurable policies to improve cluster balance, resource utilization, and scheduling efficiency.",
			Expect: true,
		},
		{
			Name:   "tofu-controller-fips",
			Input:  "A GitOps OpenTofu and Terraform controller for Flux",
			Expect: true,
		},
		{
			Name:   "opentofu-fips",
			Input:  "[OpenTofu](https://opentofu.org/) is an open-source infrastructure as code tool that allows you to declaratively manage your cloud infrastructure. OpenTofu is a fork of Terraform managed by the Linux Foundation.",
			Expect: true,
		},
		{
			Name:   "nrdot-collector-k8s",
			Input:  "New Relic .NET collector for Kubernetes monitoring",
			Expect: true,
		},
		{
			Name:   "trivy-operator",
			Input:  "The Trivy Operator automatically scans your Kubernetes workloads for security issues and generates security reports as Kubernetes Custom Resources.",
			Expect: true,
		},
		{
			Name:   "opal",
			Input:  "OPAL is an administration layer for Policy Engines",
			Expect: true,
		},
		{
			Name:   "liberica-jdk",
			Input:  "Free and open source Progressive Java Runtime for modern Java deployments",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "secrets-store-csi-driver",
			Input:  "Minimal image with Kubernetes Secrets Store CSI Driver.",
			Expect: true,
		},
		{
			Name:   "newrelic-infrastructure-bundle-fips",
			Input:  "Minimal [newrelic-infrastructure-bundle](https://github.com/newrelic/infrastructure-bundle) container image with FIPS-compliant cryptography.",
			Expect: true,
		},
		{
			Name:   "kubernetes-event-exporter-fips",
			Input:  "Minimalist [wolfi](https://github.com/wolfi-dev)-based FIPS image of [Kubernetes Event Exporter](https://github.com/resmoio/kubernetes-event-exporter). Exports Kubernetes events to various outputs to be used for observability or alerting purposes.",
			Expect: true,
		},
		{
			Name:   "aws-sigv4-proxy-fips",
			Input:  "This project signs and proxies HTTP requests with Sigv4",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "newrelic-kubernetes-fips",
			Input:  "FIPS-compliant minimal [newrelic-kubernetes](https://github.com/newrelic/nri-kubernetes) container image.",
			Expect: true,
		},
		{
			Name:   "prometheus-redis-exporter",
			Input:  "Minimalist Wolfi-based Prometheus Redis Exporter image for exporting metrics to Redis.",
			Expect: true,
		},
		{
			Name:   "grype",
			Input:  "A vulnerability scanner for container images and filesystems",
			Expect: true,
		},
		{
			Name:   "elastic-agent",
			Input:  "Elastic Agent is a unified agent for collecting, monitoring, and securing data across systems in the Elastic Stack.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "request-6335",
			Input:  "Custom Base image on top of Chainguard glibc-dynamic image.",
			Expect: true,
		},
		{
			Name:   "ntpd-rs",
			Input:  "Minimal image with [ntpd-rs](https://github.com/pendulum-project/ntpd-rs).",
			Expect: true,
		},
		{
			Name:   "karpenter",
			Input:  "Minimal image with Karpenter.",
			Expect: true,
		},
		{
			Name:   "clang",
			Input:  "[Clang](https://clang.llvm.org) is a compiler front end for the C, C++, Objective-C, and Objective-C++ programming languages, as well as the OpenMP, OpenCL, RenderScript, CUDA, SYCL, and HIP frameworks",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "openfga-fips",
			Input:  "OpenFGA is a high-performance and flexible authorization and permission engine built for developers and inspired by Google Zanzibar.",
			Expect: true,
		},
		{
			Name:   "skopeo",
			Input:  "Minimalist Wolfi-based skopeo image for interacting with container registries.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "helm-fips",
			Input:  "The Kubernetes Package Manager",
			Expect: true,
		},
		{
			Name:   "skaffold",
			Input:  "Minimal container image for running skaffold apps",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "couchdb",
			Input:  "Apache CouchDB is an open-source, document-oriented NoSQL database implemented in Erlang.",
			Expect: true,
		},
		{
			Name:   "cockroach-openssl",
			Input:  "CockroachDB is a cloud-native distributed SQL database designed to build, scale, and manage modern, data-intensive applications. The Cockroach-openssl image is the FIPS enabled equivalent of the standard Cockroach image.",
			Expect: true,
		},
		{
			Name:   "local-path-provisioner",
			Input:  "Local Path Provisioner provides a way for the Kubernetes users to utilize the local storage in each node.",
			Expect: true,
		},
		{
			Name:   "thanos-iamguarded",
			Input:  "Highly available Prometheus setup with long term storage",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "kubernetes-secret-generator",
			Input:  "Kubernetes controller for automatically generating and updating secrets",
			Expect: true,
		},
		{
			Name:   "kuberay-operator",
			Input:  "A toolkit to run Ray applications on Kubernetes",
			Expect: true,
		},
		{
			Name:   "kube-webhook-certgen",
			Input:  "Generates certificates and updates Kubernetes webhooks, integrating with Helm to simplify Kubernetes job execution.",
			Expect: true,
		},
		{
			Name:   "k8s_gateway",
			Input:  "A CoreDNS plugin to resolve all types of external Kubernetes resources",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "kyverno-notation-aws-fips",
			Input:  "Kyverno extension service for Notation and the AWS signer",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "kube-vip",
			Input:  "Kubernetes Control Plane Virtual IP and Load-Balancer",
			Expect: true,
		},
		{
			Name:   "grafana-operator",
			Input:  "A Wolfi-powered operator for Grafana that installs and manages Grafana instances, Dashboards and Datasources.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "longhorn-fips",
			Input:  "A lightweight, reliable distributed block storage system for Kubernetes.",
			Expect: true,
		},
		{
			Name:   "malcontent",
			Input:  "Enumerate binary capabilities, including malicious behaviors.",
			Expect: true,
		},
		{
			Name:   "kube-downscaler",
			Input:  "Minimal image with [kube-downscaler](https://codeberg.org/hjacobs/kube-downscaler), scale down Kubernetes deployments after work hours.",
			Expect: true,
		},
		{
			Name:   "kubeflow-katib",
			Input:  "Minimalist Kubeflow Katib Images",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "tritonserver",
			Input:  "The Triton Inference Server provides an optimized cloud and edge inferencing solution.",
			Expect: true,
		},
		{
			Name:   "temporal-ui-server",
			Input:  "Golang Server for https://github.com/temporalio/ui",
			Expect: true,
		},
		{
			Name:   "hydra",
			Input:  "Ory Hydra is a hardened, OpenID Certified OAuth 2.0 Server and OpenID Connect Provider optimized for low-latency, high throughput, and low resource consumption.",
			Expect: true,
		},
		{
			Name:   "kubescape-operator-fips",
			Input:  "Minimal Fips image of Kubescape-Operator is an in-cluster component of the Kubescape security platform that orchestrates security scanning and policy enforcement.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "nodetaint",
			Input:  "Minimal [nodetaint](https://github.com/wish/nodetaint) container image.",
			Expect: true,
		},
		{
			Name:   "steampipe",
			Input:  "Steampipe is the zero-ETL way to query APIs and services, used to expose data sources to SQL.",
			Expect: true,
		},
		{
			Name:   "prometheus-operator-iamguarded",
			Input:  "Prometheus Operator creates/configures/manages Prometheus clusters atop Kubernetes",
			Expect: true,
		},
		{
			Name:   "r-base",
			Input:  "This image contains the R programming language and environment.It can be used for statistical analysis, machine learning and data visualization.",
			Expect: true,
		},
		{
			Name:   "dragonfly-operator",
			Input:  "Kubernetes operator used to deploy and manage Dragonfly instances inside your Kubernetes clusters.",
			Expect: true,
		},
		{
			Name:   "metrics-server-iamguarded-fips",
			Input:  "Metrics Server is a Kubernetes component that collects and provides resource usage metrics (CPU, memory) for nodes and pods",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "apache-jena-fuseki",
			Input:  "SPARQL 1.1 server with a web interface, backed by the Apache Jena TDB RDF triple store.",
			Expect: true,
		},
		{
			Name:   "kubernetes-csi-external-provisioner-fips",
			Input:  "Sidecar container that watches Kubernetes PersistentVolumeClaim objects and triggers CreateVolume/DeleteVolume against a CSI endpoint",
			Expect: true,
		},
		{
			Name:   "newrelic-prometheus",
			Input:  "Minimal [newrelic-prometheus](https://github.com/newrelic/nri-prometheus) container image.",
			Expect: true,
		},
		{
			Name:   "unbound-fips",
			Input:  "Unbound is a validating, recursive, and caching DNS resolver.",
			Expect: true,
		},
		{
			Name:   "elasticsearch-iamguarded",
			Input:  "Free and Open Source, Distributed, RESTful Search Engine built with the Apache Lucene library.",
			Expect: true,
		},
		{
			Name:   "gh",
			Input:  "The GitHub CLI, or gh, is a command-line interface to GitHub for use in your terminal or your scripts.",
			Expect: true,
		},
		{
			Name:   "amazon-k8s-cni-fips",
			Input:  "Networking plugin repository for pod networking in Kubernetes using Elastic Network Interfaces on AWS",
			Expect: true,
		},
		{
			Name:   "wasmtime",
			Input:  "This image contains the `wasmtime` tool which can be used to compile or run wasm binaries.",
			Expect: true,
		},
		{
			Name:   "azuredisk-csi-fips",
			Input:  "The Azure Disk CSI driver enables the provisioning and management of Azure Disks through Kubernetes.This driver provides an interface for attaching, detaching, and managing persistent disks on Azure, helping applications achieve durable and high-performing storage.",
			Expect: true,
		},
		{
			Name:   "tritonserver-trtllm-backend",
			Input:  "The Triton Inference Server provides an optimized cloud and edge inferencing solution.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "cyberark-secrets-provider-for-k8s",
			Input:  "The CyberArk Secrets Provider for Kubernetes provides Kubernetes-based applications with access to secrets that are stored and managed in Conjur.",
			Expect: true,
		},
		{
			Name:   "cassandra-reaper",
			Input:  "Automated Repair Awesomeness for Apache Cassandra",
			Expect: true,
		},
		{
			Name:   "grafana-rollout-operator-fips",
			Input:  "Kubernetes Rollout Operator coordinates the rollout of pods between different StatefulSets within a specific namespace, and can be used to manage multi-AZ deployments",
			Expect: true,
		},
		{
			Name:   "spamcheck",
			Input:  "Spamcheck is a gRPC-based spam classification service for GitLab",
			Expect: true,
		},
		{
			Name:   "keycloak-iamguarded",
			Input:  "Minimalist Wolfi-based [Keycloak](https://www.keycloak.org/) IAMGuarded image for identity and access management.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "cass-config-builder",
			Input:  "Minimal [cass-config-builder](https://github.com/datastax/cass-config-builder) container image.",
			Expect: true,
		},
		{
			Name:   "flannel-cni-plugin",
			Input:  "A plugin designed to work in conjunction with flannel, a network fabric for containers.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "ruby-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "cyberark-secrets-provider-for-k8s-fips",
			Input:  "The CyberArk Secrets Provider for Kubernetes provides Kubernetes-based applications with access to secrets that are stored and managed in Conjur.",
			Expect: true,
		},
		{
			Name:   "python",
			Input:  "Minimal Python image based on Wolfi.",
			Expect: true,
		},
		{
			Name:   "mariadb",
			Input:  "[MariaDB](https://mariadb.org) is one of the most popular open source relational databases.",
			Expect: true,
		},
		{
			Name:   "harbor-registry",
			Input:  "A Wolf-based image for Harbor - an open-source container registry with policies and RBAC, vulnerability scans, and provides trusted image signing.",
			Expect: true,
		},
		{
			Name:   "open-webui",
			Input:  "Open WebUI is an extensible, feature-rich, and user-friendly self-hosted AI platform designed to operate entirely offline. It supports various LLM runners like Ollama and OpenAI-compatible APIs, with built-in inference engine for RAG, making it a powerful AI deployment solution.",
			Expect: true,
		},
		{
			Name:   "request-5519",
			Input:  "Custom Kubernetes ingress controller implementation for HAProxy",
			Expect: true,
		},
		{
			Name:   "wildfly",
			Input:  "WildFly is a lightweight and open-source application server designed for deploying and running Java applications.",
			Expect: true,
		},
		{
			Name:   "haproxy",
			Input:  "A minimal [haproxy](https://www.haproxy.org/) base image rebuilt every night from source.",
			Expect: true,
		},
		{
			Name:   "request-6766",
			Input:  "A custom build of ingress-nginx-controller-fips with image config changes",
			Expect: true,
		},
		{
			Name:   "akhq",
			Input:  "Kafka GUI for Apache Kafka to manage topics, topics data, consumers group, schema registry",
			Expect: true,
		},
		{
			Name:   "copybara",
			Input:  "Copybara is a tool used for transforming and moving code between repositories, enabling code synchronization workflows between different version control systems.",
			Expect: true,
		},
		{
			Name:   "mysql-client",
			Input:  "A simple SQL shell with input line editing capabilities, to interact with MySQL",
			Expect: true,
		},
		{
			Name:   "kubernetes-dashboard-web",
			Input:  "Module containing web application written in Angular and Go server with some web-related logic.",
			Expect: true,
		},
		{
			Name:   "fluentd",
			Input:  "[Fluentd](https://www.fluentd.org/): Unified Logging Layer (project under CNCF)",
			Expect: true,
		},
		{
			Name:   "blob-csi",
			Input:  "This driver allows Kubernetes to access Azure Storage via azure-storage-fuse & NFSv3.",
			Expect: true,
		},
		{
			Name:   "duckdb-fips",
			Input:  "DuckDB is an analytical in-process SQL database management system",
			Expect: true,
		},
		{
			Name:   "prometheus-mysqld-exporter-iamguarded",
			Input:  "Prometheus exporter for MySQL server metrics.",
			Expect: true,
		},
		{
			Name:   "hugo",
			Input:  "This is a minimal [Hugo](https://gohugo.io/) image.",
			Expect: true,
		},
		{
			Name:   "newrelic-kubernetes",
			Input:  "Minimal [newrelic-kubernetes](https://github.com/newrelic/nri-kubernetes) container image.",
			Expect: true,
		},
		{
			Name:   "jenkins-inbound-agent",
			Input:  "The Jenkins Inbound Agent container image is designed to run as a Jenkins agent that connects to a Jenkins controller via inbound (JNLP) connection method, used in Jenkins Kubernetes setups.",
			Expect: true,
		},
		{
			Name:   "gitlab",
			Input:  "GitLab is a complete DevOps platform that provides source code management, CI/CD automation, and collaboration tools in a single application for the entire software development lifecycle.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "db-operator",
			Input:  "The DB Operator creates databases and make them available in the cluster via Custom Resource.",
			Expect: true,
		},
		{
			Name:   "minio-object-browser-iamguarded",
			Input:  "MinIO Console is a library that provides a management and browser UI overlay for the MinIO Server",
			Expect: true,
		},
		{
			Name:   "kibana",
			Input:  "Your window into the Elastic Stack.",
			Expect: true,
		},
		{
			Name:   "nfs-subdir-external-provisioner",
			Input:  "Dynamic sub-dir volume provisioner on a remote NFS server.",
			Expect: true,
		},
		{
			Name:   "eks-distro-fips",
			Input:  "An open-source distribution of Kubernetes from AWS",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "nginx-prometheus-exporter-iamguarded-fips",
			Input:  "Nginx prometheus exporter",
			Expect: true,
		},
		{
			Name:   "x509-certificate-exporter-fips",
			Input:  "A Prometheus exporter to monitor x509 certificates expiration in Kubernetes clusters or standalone",
			Expect: true,
		},
		{
			Name:   "grafana-fips",
			Input:  "A minimal wolfi-based image for grafana, which is an open-source monitoring and observability application",
			Expect: true,
		},
		{
			Name:   "azuredisk-csi",
			Input:  "The Azure Disk CSI driver enables the provisioning and management of Azure Disks through Kubernetes",
			Expect: true,
		},
		{
			Name:   "jellyfin",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "aws-node-termination-handler-fips",
			Input:  "Gracefully handle EC2 instance shutdown within Kubernetes",
			Expect: true,
		},
		{
			Name:   "prometheus-pgbouncer-exporter",
			Input:  "A Prometheus exporter that collects and exposes metrics from PgBouncer, a lightweight connection pooler for PostgreSQL",
			Expect: true,
		},
		{
			Name:   "keycloak-config-cli-iamguarded",
			Input:  "Import YAML/JSON-formatted configuration files into Keycloak - Configuration as Code for Keycloak.",
			Expect: true,
		},
		{
			Name:   "request-4635",
			Input:  "Apache Airflow offers a platform to author, schedule, and monitor workflows programmatically. This image is a minimal, slimmed-down version of the official Apache Airflow with only core components. This image comes with pip out of the box for extending the image, but does not guarantee the image will still be FIPS compliant after installing additional components.",
			Expect: true,
		},
		{
			Name:   "kyverno-policy-reporter-plugin-trivy-fips",
			Input:  "This Plugin for Policy Reporter brings additional Trivy specific information to the Policy Reporter UI",
			Expect: true,
		},
		{
			Name:   "aws-ebs-csi-driver",
			Input:  "Minimal images for [aws-ebs-csi-driver](https://aws.amazon.com/ebs/).",
			Expect: true,
		},
		{
			Name:   "datadog-operator",
			Input:  "Kubernetes Operator for Datadog Resources",
			Expect: true,
		},
		{
			Name:   "request-5947",
			Input:  "A custom build of dcgm-exporter without setcap SYS_ADMIN",
			Expect: true,
		},
		{
			Name:   "yunikorn-scheduler-fips",
			Input:  "Apache YuniKorn K8shim",
			Expect: true,
		},
		{
			Name:   "drupal",
			Input:  "Drupal is an open source content management platform supporting a variety of websites ranging from personal weblogs to large community-driven websites.",
			Expect: true,
		},
		{
			Name:   "promitor-agent-scraper",
			Input:  "Promitor is an Azure Monitor scraper which makes the metrics available through a scraping endpoint for Prometheus or push to a StatsD server",
			Expect: true,
		},
		{
			Name:   "memcached-iamguarded-fips",
			Input:  "[Memcached](https://memcached.org/) is an in-memory key-value store for small chunks of arbitrary data (strings, objects) from results of database calls, API calls, or page rendering.",
			Expect: true,
		},
		{
			Name:   "kaniko-warmer-fips",
			Input:  "Build Container Images In Kubernetes",
			Expect: true,
		},
		{
			Name:   "consul-k8s",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "victoria-metrics",
			Input:  "VictoriaMetrics standalone image is a fast, cost-effective and scalable monitoring solution and time series database",
			Expect: true,
		},
		{
			Name:   "rabbitmq-default-user-credential-updater-fips",
			Input:  "Minimal FIPS compliant image with [default-user-credential-updater](https://github.com/rabbitmq/default-user-credential-updater)",
			Expect: true,
		},
		{
			Name:   "sealed-secrets-kubeseal-iamguarded",
			Input:  "A Kubernetes tool for one-way encrypted Secrets, enabling safe GitOps-friendly secret management.",
			Expect: true,
		},
		{
			Name:   "telegraf",
			Input:  "Minimal image with Telegraf agent for collecting, processing, aggregating, and writing metrics.",
			Expect: true,
		},
		{
			Name:   "apache-tika",
			Input:  "Apache Tika extracts metadata, text, and language from documents, enabling content analysis and indexing.",
			Expect: true,
		},
		{
			Name:   "images/grafana-alloy-operator-fips",
			Input:  "The Alloy Operator is a Kubernetes Operator that manages the lifecycle of Grafana Alloy instances",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "tomcat",
			Input:  "[Tomcat](https://tomcat.apache.org/) is a free and open-source implementation of the Jakarta Servlet, Jakarta Expression Language, and WebSocket technologies.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "cilium-certgen",
			Input:  "A convenience tool to generate and store certificates for Hubble Relay mTLS",
			Expect: true,
		},
		{
			Name:   "quic-go-fips",
			Input:  "A production-ready QUIC implementation in pure Go",
			Expect: true,
		},
		{
			Name:   "go-fips",
			Input:  "Container image for building Go applications with FIPS",
			Expect: true,
		},
		{
			Name:   "volsync-fips",
			Input:  "Asynchronous data replication for Kubernetes volumes",
			Expect: true,
		},
		{
			Name:   "mysql-fips",
			Input:  "MySQL is a widely used open-source relational database management system.",
			Expect: true,
		},
		{
			Name:   "melange",
			Input:  "Container image for running [melange](https://github.com/chainguard-dev/melange) workflows to build APK packages.",
			Expect: true,
		},
		{
			Name:   "rails",
			Input:  "Ruby on Rails (often just called \"Rails\") is a web-application framework that includes everything needed to create database-backed web applications according to the Model-View-Controller (MVC) pattern.",
			Expect: true,
		},
		{
			Name:   "amazon-corretto-jre-fips",
			Input:  "Amazon Corretto is a no-cost, multi-platform, production-ready distribution of corresponding version of OpenJDK",
			Expect: true,
		},
		{
			Name:   "apache-hop-fips",
			Input:  "Data orchestration and engineering platform for managing ETL/ELT workflows.",
			Expect: true,
		},
		{
			Name:   "local-volume-node-cleanup",
			Input:  "The local volume node cleanup controller removes PersistentVolumes and PersistentVolumeClaims that reference deleted Nodes.",
			Expect: true,
		},
		{
			Name:   "ztunnel",
			Input:  "The ztunnel component of ambient mesh",
			Expect: true,
		},
		{
			Name:   "images/valkey-bundle-fips",
			Input:  "Valkey bundle with pre-installed modules for extended functionality.",
			Expect: true,
		},
		{
			Name:   "pgbouncer-iamguarded-fips",
			Input:  "This image contains the CLI for the [pgbouncer](https://www.pgbouncer.org/) connection pooler for PostgreSQL with IAMGuarded support.",
			Expect: true,
		},
		{
			Name:   "mesosphere-vsphere-csi-driver",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "cert-manager-iamguarded-fips",
			Input:  "[cert-manager](https://cert-manager.io) is a tool for provisioning and managing TLS certificates in Kubernetes.",
			Expect: true,
		},
		{
			Name:   "splunk-otel-collector",
			Input:  "Splunk OpenTelemetry Collector is a distribution of the OpenTelemetry Collector. It provides a unified way to receive, process, and export metric, trace, and log data for Splunk Observability Cloud",
			Expect: true,
		},
		{
			Name:   "backup-restore-operator",
			Input:  "A Backup and Restore Operator that provides the ability to back up and restore the Rancher application running on any Kubernetes cluster.",
			Expect: true,
		},
		{
			Name:   "nacos",
			Input:  "Dynamic service discovery, configuration and service management platform for building AI cloud native applications.",
			Expect: true,
		},
		{
			Name:   "vault-csi-provider-fips",
			Input:  "HashiCorp Vault Provider for Secret Store CSI Driver",
			Expect: true,
		},
		{
			Name:   "metrics-server",
			Input:  "Metrics Server is a Kubernetes component that collects and provides resource usage metrics (CPU, memory) for nodes and pods.",
			Expect: true,
		},
		{
			Name:   "bats-fips",
			Input:  "Bats provides a simple way to verify that the UNIX programs you write behave as expected.",
			Expect: true,
		},
		{
			Name:   "openresty",
			Input:  "OpenResty is a high Performance Web Platform Based on Nginx and LuaJIT.",
			Expect: true,
		},
		{
			Name:   "rancher-webhook-fips",
			Input:  "Rancher Webhook automates tasks and integrates external systems in response to events within a Rancher-managed Kubernetes environment.",
			Expect: true,
		},
		{
			Name:   "netcat",
			Input:  "Minimal image for Debian port of OpenBSD's netcat.",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "ip-masq-agent",
			Input:  "Minimal image to manage IP masquerading on Kubernetes nodes",
			Expect: true,
		},
		{
			Name:   "external-dns-iamguarded-fips",
			Input:  "Minimal image to configure external DNS servers (AWS Route53, Google CloudDNS and others) for Kubernetes Ingresses and Services",
			Expect: true,
		},
		{
			Name:   "null",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "kubernetes-event-exporter-iamguarded-fips",
			Input:  "IAMGuarded compatible image of [Kubernetes Event Exporter](https://github.com/resmoio/kubernetes-event-exporter) for exporting Kubernetes events to various outputs to be used for observability or alerting purposes.",
			Expect: true,
		},
		{
			Name:   "az-iamguarded",
			Input:  "Azure CLI (IAM Guarded)",
			Expect: true,
		},
		{
			Name:   "clickhouse-fips",
			Input:  "Clickhouse is the fastest and most resource efficient open-source database for real-time apps and analytics, built with FIPS 140-3 compliant cryptographic modules.",
			Expect: true,
		},
		{
			Name:   "wordpress",
			Input:  "Minimalist Wolfi-based WordPress images.",
			Expect: true,
		},
		{
			Name:   "cluster-api",
			Input:  "Home for Cluster API, a subproject of sig-cluster-lifecycle",
			Expect: true,
		},
		{
			Name:   "kubernetes-csi-external-provisioner",
			Input:  "Sidecar container that watches Kubernetes PersistentVolumeClaim objects and triggers CreateVolume/DeleteVolume against a CSI endpoint",
			Expect: true,
		},
		{
			Name:   "cluster-api-ipam-provider-in-cluster",
			Input:  "An IPAM provider for Cluster API that manages pools of IP addresses using Kubernetes resources.",
			Expect: true,
		},
		{
			Name:   "prometheus-podman-exporter",
			Input:  "Prometheus exporter for podman environments exposing containers, pods, images, volumes and networks information.",
			Expect: true,
		},
		{
			Name:   "istio",
			Input:  "[Istio](https://istio.io) is a service mesh that extends Kubernetes to provide traffic management, telemetry, security, and policy for complex deployments.",
			Expect: true,
		},
		{
			Name:   "spiffe-helper-fips",
			Input:  "A secure, minimal container image for the SPIFFE Helper utility that automates X.509 SVID certificate rotation for services that can't natively fetch X.509-SVIDs.",
			Expect: true,
		},
		{
			Name:   "kubernetes-dns-node-cache",
			Input:  "Minimal image that acts as a drop-in replacement for the [NodeLocal DNSCache](https://github.com/kubernetes/dns) image.",
			Expect: true,
		},
		{
			Name:   "valkey-sentinel-iamguarded-fips",
			Input:  "Valkey is an open source, in-memory data store used by millions of developers as a cache, vector database, document database, streaming engine, and message broker.",
			Expect: true,
		},
		{
			Name:   "task-fips",
			Input:  "Task is a task runner and build tool that aims to be simpler and easier to use than GNU Make.",
			Expect: true,
		},
		{
			Name:   "prometheus-alertmanager-iamguarded-fips",
			Input:  "The Alertmanager handles alerts sent by client applications such as the Prometheus server",
			Expect: true,
		},
		{
			Name:   "kubernetes-event-exporter-iamguarded",
			Input:  "IAMGuarded compatible image of [Kubernetes Event Exporter](https://github.com/resmoio/kubernetes-event-exporter) for exporting Kubernetes events to various outputs to be used for observability or alerting purposes.",
			Expect: true,
		},
		{
			Name:   "rclone",
			Input:  "Rclone syncs files and directories to and from different cloud storage providers.",
			Expect: true,
		},
		{
			Name:   "langchain",
			Input:  "LangChainPython library image.",
			Expect: true,
		},
		{
			Name:   "request-6130",
			Input:  "Minimal image with JuiceFS. JuiceFS is a high-performance POSIX file system, particularly designed for the cloud-native environment.",
			Expect: true,
		},
		{
			Name:   "elastic-agent-fips",
			Input:  "Elastic Agent is a unified agent for collecting, monitoring, and securing data across systems in the Elastic Stack.",
			Expect: true,
		},
		{
			Name:   "seata-server",
			Input:  "Seata is a high-performance, easy-to-use distributed transaction solution designed for microservices architecture.",
			Expect: true,
		},
		{
			Name:   "tileserver-gl",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "git-iamguarded-fips",
			Input:  "A minimal Git image for use with Iamguarded charts.",
			Expect: true,
		},
		{
			Name:   "request-7270",
			Input:  "Minimal image with [opentelemetry-collector](https://github.com/open-telemetry/opentelemetry-collector)",
			Expect: true,
		},
		{
			Name:   "flux",
			Input:  "`flux` cli to interact with the [Flux](https://fluxcd.io/) gitops toolkit components in a running cluster.",
			Expect: true,
		},
		{
			Name:   "confluent-kafka",
			Input:  "A Wolfi-based container image for the Community Edition of Confluent Kafka (cp-kafka), which extends Apache Kafka with additional features.",
			Expect: true,
		},
		{
			Name:   "spicedb-operator",
			Input:  "SpiceDB Operator is a Kubernetes operator for managing SpiceDB clusters, providing automated deployment, scaling, and management of SpiceDB instances.",
			Expect: true,
		},
		{
			Name:   "shadowsocks-rust",
			Input:  "Shadowsocks-rust is a Rust implementation of the Shadowsocks protocol, aimed at ensuring secure and private internet access by encrypting connections and circumventing internet restrictions.",
			Expect: true,
		},
		{
			Name:   "flink",
			Input:  "Apache Flink is an open source stream processing framework with powerful stream- and batch-processing capabilities.",
			Expect: true,
		},
		{
			Name:   "open-liberty",
			Input:  "Open Liberty is a highly composable, fast to start, dynamic application server runtime environment.",
			Expect: true,
		},
		{
			Name:   "oauth2-proxy-iamguarded",
			Input:  "[OAuth2 Proxy](https://oauth2-proxy.github.io/oauth2-proxy/) is a reverse proxy that provides authentication with Google, Azure, OpenID Connect and many more identity providers.",
			Expect: true,
		},
		{
			Name:   "conda",
			Input:  "This image contains the CLI for the [Conda](https://docs.conda.io/en/latest/) programming environment.",
			Expect: true,
		},
		{
			Name:   "lvm-driver-fips",
			Input:  "Dynamically provision Stateful Persistent Node-Local Volumes & Filesystems for Kubernetes that is integrated with a backend LVM2 data storage stack.",
			Expect: true,
		},
		{
			Name:   "kubernetes-dashboard-auth",
			Input:  "Go module handling authentication to the Kubernetes API.",
			Expect: true,
		},
		{
			Name:   "prometheus-admission-webhook",
			Input:  "Admission webhook for Prometheus",
			Expect: true,
		},
		{
			Name:   "aws-cli-iamguarded-fips",
			Input:  "Minimal FIPS-compliant IAMGuarded image for the [AWS CLI](https://aws.amazon.com/cli/)",
			Expect: true,
		},
		{
			Name:   "kaniko-fips",
			Input:  "Build Container Images In Kubernetes",
			Expect: true,
		},
		{
			Name:   "amazon-corretto-jdk-fips",
			Input:  "Amazon Corretto is a no-cost, multi-platform, production-ready distribution of corresponding version of OpenJDK",
			Expect: true,
		},
		{
			Name:   "gpu-operator",
			Input:  "The NVIDIA GPU Operator bootstraps, configures, and manages GPUs in Kubernetes.",
			Expect: true,
		},
		{
			Name:   "cloud-provider-aws",
			Input:  "Cloud provider for AWS",
			Expect: true,
		},
		{
			Name:   "omni",
			Input:  "Omni managing Kubernetes clusters across bare metal, VMs, and cloud",
			Expect: true,
		},
		{
			Name:   "heartbeat-fips",
			Input:  "Heartbeat periodically checks the status of your services and determine whether they are available.",
			Expect: true,
		},
		{
			Name:   "minio-operator",
			Input:  "MinIO Operator - Operator for MinIO on Kubernetes",
			Expect: true,
		},
		{
			Name:   "kubo-fips",
			Input:  "FIPS-compliant, minimalist Wolfi-based container image for the IPFS Kubo (go-ipfs) node.",
			Expect: true,
		},
		{
			Name:   "oauth2-proxy",
			Input:  "[OAuth2 Proxy](https://oauth2-proxy.github.io/oauth2-proxy/) is a reverse proxy that provides authentication with Google, Azure, OpenID Connect and many more identity providers.",
			Expect: true,
		},
		{
			Name:   "cert-manager-istio-csr-fips",
			Input:  "istio-csr is an agent that allows for Istio workload and control plane components to be secured using cert-manager.",
			Expect: true,
		},
		{
			Name:   "atmoz-sftp-fips",
			Input:  "Various scripts from https://github.com/atmoz/sftp to help run SFTP in a container",
			Expect: true,
		},
		{
			Name:   "eclipse-mosquitto",
			Input:  "Eclipse Mosquitto - An open source MQTT broker",
			Expect: true,
		},
		{
			Name:   "gitlab-operator",
			Input:  "Kubernetes Operator for GitLab Server",
			Expect: true,
		},
		{
			Name:   "logstash",
			Input:  "Logstash is a server-side data processing pipeline that ingests data from multiple sources, transforms it, and sends it to your chosen destination",
			Expect: true,
		},
		{
			Name:   "request-5278",
			Input:  "A custom image based on the Chainguard image container image for Confluent Kafka (community edition), with some modifications per customer requirements.",
			Expect: true,
		},
		{
			Name:   "rancher-agent",
			Input:  "The Rancher Agent is responsible for communicating between the managed Kubernetes clusters and the Rancher server",
			Expect: true,
		},
		{
			Name:   "prometheus-mysqld-exporter-iamguarded-fips",
			Input:  "Prometheus exporter for MySQL server metrics.",
			Expect: true,
		},
		{
			Name:   "valkey-sentinel-iamguarded",
			Input:  "Valkey is an open source, in-memory data store used by millions of developers as a cache, vector database, document database, streaming engine, and message broker.",
			Expect: true,
		},
		{
			Name:   "kubernetes-csi-external-snapshotter",
			Input:  "Sidecar container that watches Kubernetes Snapshot CRD objects and triggers CreateSnapshot/DeleteSnapshot against a CSI endpoint. This container image includes the snapshotter, snapshot-controller and snapshot-validation-webhook.",
			Expect: true,
		},
		{
			Name:   "langfuse-fips",
			Input:  "Langfuse is an open-source observability and analytics platform for LLM applications.",
			Expect: true,
		},
		{
			Name:   "debezium-connect",
			Input:  "Kafka Connect image with all Debezium connectors, and part of the Debezium platform",
			Expect: true,
		},
		{
			Name:   "rqlite",
			Input:  "Minimal image with rqlite.",
			Expect: true,
		},
		{
			Name:   "kepler",
			Input:  "Kepler (Kubernetes-based Efficient Power Level Exporter) is a Prometheus exporter that measures energy consumption metrics at the container, pod, and node levels in Kubernetes clusters.",
			Expect: true,
		},
		{
			Name:   "argocd-extension-installer",
			Input:  "Install Argo CD extensions using init-containers",
			Expect: true,
		},
		{
			Name:   "filebeat",
			Input:  "[filebeat](https://github.com/elastic/beats/tree/main/filebeat) Tails and ships log files",
			Expect: true,
		},
		{
			Name:   "redis-iamguarded-fips",
			Input:  "[Redis](https://github.com/redis/redis) Redis is an in-memory data structure store, used as a database, cache, and message broker.",
			Expect: true,
		},
		{
			Name:   "elixir",
			Input:  "Elixir is a dynamic, functional language for building scalable and maintainable applications.",
			Expect: true,
		},
		{
			Name:   "boring-registry",
			Input:  "Minimal image with the `boring-registry` [server application](https://github.com/TierMobility/boring-registry).",
			Expect: true,
		},
		{
			Name:   "jaeger",
			Input:  "CNCF Jaeger, a Distributed Tracing Platform",
			Expect: true,
		},
		{
			Name:   "librechat",
			Input:  "An open-source AI chat application.",
			Expect: true,
		},
		{
			Name:   "verticadb-operator-fips",
			Input:  "The VerticaDB operator automates tasks and monitors the state of your Vertica on Kubernetes deployments.",
			Expect: true,
		},
		{
			Name:   "timescaledb-compat",
			Input:  "A time-series database for high-performance real-time analytics packaged as a Postgres extension",
			Expect: true,
		},
		{
			Name:   "grafana-mimir",
			Input:  "A minimal Wolfi-based image for Grafana Mimir, providing horizontally scalable, highly available, multi-tenant, long-term storage for Prometheus.",
			Expect: true,
		},
		{
			Name:   "thanos-receive-controller-fips",
			Input:  "Kubernetes controller to automatically configure Thanos receive hashrings",
			Expect: true,
		},
		{
			Name:   "cert-exporter",
			Input:  "A minimal, wolfi-based image for cert-exporter: an application that exports certificate expiration metrics from disk, Kubernetes, and AWS Secrets Manager to Prometheus.",
			Expect: true,
		},
		{
			Name:   "camunda-keycloak-fips",
			Input:  "Minimalist Wolfi-based [Camunda Keycloak](https://github.com/camunda/keycloak/) image for identity and access management.",
			Expect: true,
		},
		{
			Name:   "dnsdist",
			Input:  "dnsdist is a highly DNS-, DoS- and abuse-aware loadbalancer",
			Expect: true,
		},
		{
			Name:   "node-feature-discovery",
			Input:  "A minimal wolfi-based image for node-feature-discovery, Node feature discovery for Kubernetes",
			Expect: true,
		},
		{
			Name:   "redis",
			Input:  "Chainguard image for [Redis](https://github.com/redis/redis), an in-memory database that persists on disk. Redis is a key-value store, supporting an array of different values, including Strings, Lists, Sets, Sorted Sets, Hashes, Streams, HyperLogLogs, and Bitmaps.",
			Expect: true,
		},
		{
			Name:   "adoptium-jre",
			Input:  "Minimalist Wolfi-based Java JRE image using [Adoptium](https://adoptium.net/). Used for running Java applications.",
			Expect: true,
		},
		{
			Name:   "mdbook",
			Input:  "Minimal image that contains [mdbook](https://rust-lang.github.io/mdBook/).",
			Expect: true,
		},
		{
			Name:   "nrdot-collector-k8s-fips",
			Input:  "New Relic .NET collector for Kubernetes monitoring",
			Expect: true,
		},
		{
			Name:   "docker-compose-fips",
			Input:  "minimal docker-compose image with docker-compose binary",
			Expect: true,
		},
		{
			Name:   "nginx-iamguarded-fips",
			Input:  "Minimal Wolfi-based nginx HTTP, reverse proxy, mail proxy, and a generic TCP/UDP proxy server",
			Expect: true,
		},
		{
			Name:   "kubernetes-csi-driver-nfs",
			Input:  "This driver allows Kubernetes to access an NFS server running on a Linux node.",
			Expect: true,
		},
		{
			Name:   "cert-manager-csi-driver",
			Input:  "A Kubernetes CSI driver that automatically mounts signed certificates to Pods using ephemeral volumes",
			Expect: true,
		},
		{
			Name:   "metricbeat",
			Input:  "Metricbeat fetches a set of metrics on a predefined interval from the operating system and services such as Apache web server, Redis, and more and ships them to Elasticsearch or Logstash.",
			Expect: true,
		},
		{
			Name:   "redisinsight",
			Input:  "Redis GUI by Redis",
			Expect: true,
		},
		{
			Name:   "monstache",
			Input:  "A go daemon that syncs mongodb to elasticsearch in realtime.",
			Expect: true,
		},
		{
			Name:   "kubernetes-dashboard-metrics-scraper",
			Input:  "Module containing the Kubernetes metrics scraper module of the Kubernetes dashboard application",
			Expect: true,
		},
		{
			Name:   "ingress-nginx-controller-iamguarded",
			Input:  " Ingress-NGINX Controller for Kubernetes",
			Expect: true,
		},
		{
			Name:   "headlamp-plugin-flux",
			Input:  "Headlamp plugin to visualize and manage Flux GitOps resources in Kubernetes clusters.",
			Expect: true,
		},
		{
			Name:   "image-factory-fips",
			Input:  "Image Factory is a flexible artifact-build service for Talos Linux, offering automated creation of ISOs, disks, cloud images, and installers and ensuring reproductible builds.",
			Expect: true,
		},
		{
			Name:   "semgrep",
			Input:  "CLI for the [Semgrep](https://semgrep.dev) static analysis tool. Semgrep is a lightweight static analysis tool for many languages. It finds bug variants with patterns that look like source code.",
			Expect: true,
		},
		{
			Name:   "ceph",
			Input:  "Ceph is a distributed object, block, and file storage platform",
			Expect: true,
		},
		{
			Name:   "dive",
			Input:  "Minimal [dive](https://github.com/wagoodman/dive) container image.",
			Expect: true,
		},
		{
			Name:   "newrelic-nri-statsd",
			Input:  "The StatsD integration lets you easily get StatsD data into New Relic",
			Expect: true,
		},
		{
			Name:   "auditbeat-fips",
			Input:  "Auditbeat is a lightweight shipper that you can install on your servers to audit the activities of users and processes on your systems.",
			Expect: true,
		},
		{
			Name:   "mariadb-operator",
			Input:  "Mariadb-operator is a Kubernetes operator for managing MariaDB databases. It automates the deployment, scaling, and management of MariaDB instances in Kubernetes clusters, providing declarative configuration and lifecycle management.",
			Expect: true,
		},
		{
			Name:   "kubernetes-ingress-defaultbackend",
			Input:  "Minimal image that acts as a drop-in replacement for the `registry.k8s.io/defaultbackend` image. Used in some ingresses like https://github.com/kubernetes/ingress-gce and https://github.com/kubernetes/ingress-nginx",
			Expect: true,
		},
		{
			Name:   "dart",
			Input:  "Container image for dart programming language",
			Expect: true,
		},
		{
			Name:   "kiali-operator-fips",
			Input:  "A FIPS compliant Kiali Operator that manages the lifecycle of Kiali in Kubernetes environments with Istio service mesh integration",
			Expect: true,
		},
		{
			Name:   "spiffe-helper",
			Input:  "A secure, minimal container image for the SPIFFE Helper utility that automates X.509 SVID certificate rotation for services that can't natively fetch X.509-SVIDs.",
			Expect: true,
		},
		{
			Name:   "jetty",
			Input:  "Eclipse Jetty is a lightweight, highly scalable, Java-based web server and Servlet engine",
			Expect: true,
		},
		{
			Name:   "seaweedfs-operator",
			Input:  "seaweedfs kubernetes operator",
			Expect: true,
		},
		{
			Name:   "gcp-compute-persistent-disk-csi-driver",
			Input:  "The Google Compute Engine Persistent Disk (GCE PD) Container Storage Interface (CSI) Storage Plugin.",
			Expect: true,
		},
		{
			Name:   "request-5166",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "terragrunt",
			Input:  "A cloud infrastructure orchestration tool that supports OpenTofu/Terraform.",
			Expect: true,
		},
		{
			Name:   "k3s-openssl",
			Input:  "Minimal image with [K3s](https://k3s.io/) and our FIPS OpenSSL provider.",
			Expect: true,
		},
		{
			Name:   "kafka-bridge",
			Input:  "HTTP-based bridge for Apache Kafka using Vert.x framework",
			Expect: true,
		},
		{
			Name:   "external-dns-iamguarded",
			Input:  "Minimal image to configure external DNS servers (AWS Route53, Google CloudDNS and others) for Kubernetes Ingresses and Services",
			Expect: true,
		},
		{
			Name:   "knative-eventing",
			Input:  "Event-driven application platform for Kubernetes",
			Expect: true,
		},
		{
			Name:   "prometheus-pushgateway-iamguarded-fips",
			Input:  "Push acceptor for ephemeral and batch jobs.",
			Expect: true,
		},
		{
			Name:   "tekton-cli-fips",
			Input:  "The Tekton Pipelines CLI project provides a command-line interface (CLI) for interacting with Tekton, an open-source framework for Continuous Integration and Delivery (CI/CD) systems.",
			Expect: true,
		},
		{
			Name:   "local-volume-node-cleanup-fips",
			Input:  "The local volume node cleanup controller removes PersistentVolumes and PersistentVolumeClaims that reference deleted Nodes.",
			Expect: true,
		},
		{
			Name:   "dcgm-exporter",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "amazon-cloudwatch-agent-operator-fips",
			Input:  "A minimal, [wolfi](https://github.com/wolfi-dev)-based FIPS compliant image of  Amazon CloudWatch Agent Operator developed to manage the CloudWatch Agent on kubernetes.",
			Expect: true,
		},
		{
			Name:   "grypedb",
			Input:  "A Wolfi-powered image for GrypeDB and Vunnel that allows you to create your own Grype DB instances.",
			Expect: true,
		},
		{
			Name:   "nemo",
			Input:  "NVIDIA NeMo Framework is an end-to-end, cloud-native framework to build, customize, and deploy generative AI models anywhere.",
			Expect: true,
		},
		{
			Name:   "buildkit-fips",
			Input:  "Buildkit is a concurrent, cache-efficient, and Dockerfile-agnostic builder toolkit.",
			Expect: true,
		},
		{
			Name:   "prometheus-podman-exporter-fips",
			Input:  "Prometheus exporter for podman environments exposing containers, pods, images, volumes and networks information.",
			Expect: true,
		},
		{
			Name:   "azurefile-csi-fips",
			Input:  "This driver allows Kubernetes to access Azure File volume using smb and nfs protocols, csi plugin name: file.csi.azure.com.",
			Expect: true,
		},
		{
			Name:   "kubo",
			Input:  "Minimalist Wolfi-based container image for the IPFS Kubo (go-ipfs) node.",
			Expect: true,
		},
		{
			Name:   "auditbeat",
			Input:  "Auditbeat is a lightweight shipper that you can install on your servers to audit the activities of users and processes on your systems.",
			Expect: true,
		},
		{
			Name:   "jdk",
			Input:  "Minimal Wolfi-based Java JDK image using [OpenJDK](https://openjdk.org/projects/jdk/).  Used for compiling Java applications.",
			Expect: true,
		},
		{
			Name:   "karma-fips",
			Input:  "A dashboard for managing alerts from Alertmanager",
			Expect: true,
		},
		{
			Name:   "trivy",
			Input:  "Aquasec Trivy",
			Expect: true,
		},
		{
			Name:   "trufflehog-fips",
			Input:  "TruffleHog is a tool that allows you to discover, classify, validate, and analyze leaked credentials.",
			Expect: true,
		},
		{
			Name:   "clickhouse-keeper",
			Input:  "ClickHouse Keeper is a distributed coordination service that provides a ZooKeeper-compatible API for managing ClickHouse clusters. It handles distributed consensus, configuration management, and leader election using the Raft algorithm.",
			Expect: true,
		},
		{
			Name:   "trino",
			Input:  "Fast distributed SQL query engine for big data analytics that helps you explore your data universe.",
			Expect: true,
		},
		{
			Name:   "sealed-secrets-kubeseal-iamguarded",
			Input:  "A Kubernetes tool for one-way encrypted Secrets, enabling safe GitOps-friendly secret management.",
			Expect: true,
		},
		{
			Name:   "kubewatch",
			Input:  "[kubewatch](https://github.com/robusta-dev/kubewatch) is a Kubernetes watcher that publishes notification to available collaboration hubs/notification channels. Run it in your k8s cluster, and you will get event notifications through webhooks.",
			Expect: true,
		},
		{
			Name:   "camunda-zeebe",
			Input:  "Zeebe is the process automation engine powering Camunda.",
			Expect: true,
		},
		{
			Name:   "prometheus-node-exporter",
			Input:  "Minimalist Wolfi-based Prometheus Node Exporter image for exporting node metrics.",
			Expect: true,
		},
		{
			Name:   "thingsboard",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "kafka-proxy",
			Input:  "Proxy connections to Kafka cluster. Connect through SOCKS Proxy, HTTP Proxy or to cluster running in Kubernetes.",
			Expect: true,
		},
		{
			Name:   "gatekeeper-fips",
			Input:  "Minimal FIPS-compliant [Gatekeeper](https://open-policy-agent.github.io/gatekeeper) image for enforcing Kubernetes policies using Open Policy Agent",
			Expect: true,
		},
		{
			Name:   "os-shell-cassandra-iamguarded",
			Input:  "os-shell-cassandra-iamguarded image is an utilities image specifically for the Cassandra Helm chart's dynamic seed discovery feature. This image includes additional dependencies required for dynamic seed discovery and should only be used with the Cassandra Helm chart",
			Expect: true,
		},
		{
			Name:   "request-4077",
			Input:  "A minimal, [wolfi](https://github.com/wolfi-dev)-based image with cuda runtime and devel packages with python",
			Expect: true,
		},
		{
			Name:   "jaeger-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "zulu-jdk",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "kafbat-ui-fips",
			Input:  "FIPS compliant variant [Kafbat-UI](https://kafbat.io/) is an open-source Web UI for managing Apache Kafka clusters",
			Expect: true,
		},
		{
			Name:   "camunda",
			Input:  "Camunda is a scalable platform for process orchestration and automation",
			Expect: true,
		},
		{
			Name:   "request-5762",
			Input:  "OpenTofu-FIPS with FIPS variant Terraform providers, meant for a Custom Assembly package bundle",
			Expect: true,
		},
		{
			Name:   "backup-restore-operator-fips",
			Input:  "A Backup and Restore Operator that provides the ability to back up and restore the Rancher application running on any Kubernetes cluster.",
			Expect: true,
		},
		{
			Name:   "awx",
			Input:  "Chainguard image for AWX. Built on top of Ansible, AWX  is an automation controller that provides a web based interface, REST API, and task engine.",
			Expect: true,
		},
		{
			Name:   "glibc-openssl-fips",
			Input:  "The GNU C Library (glibc) is a C standard library implementation maintained by the GNU Project. This container image also contains OpenSSL, a software library for applications providing secure communications over a network.",
			Expect: true,
		},
		{
			Name:   "op-geth",
			Input:  "op-geth is an Ethereum execution client optimized for the Optimism Layer 2 scaling solution",
			Expect: true,
		},
		{
			Name:   "crossplane-aws-provider-fips",
			Input:  "Crossplane provider-aws is the infrastructure provider for Amazon Web Services (AWS).",
			Expect: true,
		},
		{
			Name:   "amazon-corretto-jre",
			Input:  "Amazon Corretto is a no-cost, multi-platform, production-ready distribution of corresponding version of OpenJDK",
			Expect: true,
		},
		{
			Name:   "crossplane-azure",
			Input:  "Crossplane lets you build a control plane with Kubernetes-style declarative and API-driven configuration and management for anything",
			Expect: true,
		},
		{
			Name:   "cfssl",
			Input:  "[CFSSL](https://cfssl.org/) is Cloudflare's PKI and TLS toolkit",
			Expect: true,
		},
		{
			Name:   "apache-apisix-iamguarded",
			Input:  "APISIX API Gateway provides rich traffic management features such as load balancing, dynamic upstream, canary release, circuit breaking, authentication, observability, and more.",
			Expect: true,
		},
		{
			Name:   "argocd",
			Input:  "[Argo CD](https://argo-cd.readthedocs.io/en/stable/) is a declarative continuous deployment tool for Kubernetes.",
			Expect: true,
		},
		{
			Name:   "mongodb-fips",
			Input:  "[MongoDB](https://www.mongodb.com/) is a document-oriented database management system. MongoDB is a popular example of a NoSQL database, and stores data in JSON-like documents.",
			Expect: true,
		},
		{
			Name:   "velero-restore-helper",
			Input:  "Backup and migrate Kubernetes applications and their persistent volumes",
			Expect: true,
		},
		{
			Name:   "kayenta-fips",
			Input:  "FIPS-compliant Automated Canary Service",
			Expect: true,
		},
		{
			Name:   "aws-gateway-controller",
			Input:  "AWS Application Networking is an implementation of the Kubernetes Gateway API.",
			Expect: true,
		},
		{
			Name:   "external-dns",
			Input:  "Minimal image to configure external DNS servers (AWS Route53, Google CloudDNS and others) for Kubernetes Ingresses and Services",
			Expect: true,
		},
		{
			Name:   "secrets-store-csi-driver-provider-azure",
			Input:  "Azure Key Vault provider for Secret Store CSI driver allows you to get secret contents stored in Azure Key Vault instance and use the Secret Store CSI driver interface to mount them into Kubernetes pods.",
			Expect: true,
		},
		{
			Name:   "huggingface-pytorch-inference",
			Input:  "General-purpose PyTorch inference Deep Learning Container for inferencing on Sagemaker instances",
			Expect: true,
		},
		{
			Name:   "pixi",
			Input:  "Cross-platform, multi-language package manager and workflow tool built on the foundation of the conda ecosystem.",
			Expect: true,
		},
		{
			Name:   "crane-fips",
			Input:  "Minimalist Wolfi-based crane-fips image to interact with container registries. Crane is used for inspecting and manipulating container images, allowing you to view manifests, verify image layers, and check cryptographic signatures.",
			Expect: true,
		},
		{
			Name:   "jupyterhub-k8s-network-tools",
			Input:  "Network diagnostic tools for use within a JupyterHub Kubernetes cluster",
			Expect: true,
		},
		{
			Name:   "trivy-operator-fips",
			Input:  "A FIPS-compliant container image for Trivy Operator, providing automated security scanning for Kubernetes workloads with enhanced cryptographic compliance.",
			Expect: true,
		},
		{
			Name:   "grafana-alloy",
			Input:  "OpenTelemetry Collector distribution with programmable pipelines",
			Expect: true,
		},
		{
			Name:   "tekton-chains",
			Input:  "Tekton Chains is a Kubernetes Custom Resource Definition (CRD) controller that allows you to manage your supply chain security in Tekton.",
			Expect: true,
		},
		{
			Name:   "kyverno-policy-reporter",
			Input:  "Monitoring and Observability Tool for the [PolicyReport CRD](https://kyverno.github.io/policy-reporter/) with an optional UI.",
			Expect: true,
		},
		{
			Name:   "azure-functions-node",
			Input:  "Azure Functions is a managed platform-as-a-service provider which offers scalable and serverless hosting for code projects. It extends the Azure platform with the capability to implement code triggered by many events occurring in Azure, on-premises or other 3rd party service.",
			Expect: true,
		},
		{
			Name:   "request-3495",
			Input:  "A minimal, [wolfi](https://github.com/wolfi-dev)-based image for vllm, a high-throughput and memory-efficient inference and serving engine for LLMs",
			Expect: true,
		},
		{
			Name:   "liberica-jre-fips",
			Input:  "Free and open source Progressive Java Runtime for modern Java deployments",
			Expect: true,
		},
		{
			Name:   "headlamp",
			Input:  "Headlamp is an easy-to-use and extensible Kubernetes web UI designed for developers and cluster operators.",
			Expect: true,
		},
		{
			Name:   "opensearch",
			Input:  "Minimal image with Opensearch.",
			Expect: true,
		},
		{
			Name:   "k8s-sidecar",
			Input:  "Minimal image with the k8s-sidecar app.",
			Expect: true,
		},
		{
			Name:   "crossplane-function-environment-configs-fips",
			Input:  "Crossplane function that manages environment-specific configurations for resources in compositions with FIPS 140-3 compliance",
			Expect: true,
		},
		{
			Name:   "crane",
			Input:  "Minimalist Wolfi-based crane image to interact with container registries. Crane is used for inspecting and manipulating container images, allowing you to view manifests, verify image layers, and check cryptographic signatures.",
			Expect: true,
		},
		{
			Name:   "spdx-tools",
			Input:  "A command-line utility for creating, converting, comparing, and validating SPDX documents across multiple formats.",
			Expect: true,
		},
		{
			Name:   "zulu-jre",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "kubernetes-csi-external-resizer-fips",
			Input:  "Minimal image with [kubernetes-csi/external-resizer](https://github.com/kubernetes-csi/external-resizer).",
			Expect: true,
		},
		{
			Name:   "pgpool2-iamguarded-fips",
			Input:  "open-source middleware that operates between PostgreSQL servers and clients, providing features such as connection pooling, load balancing, and replication to enhance database performance and availability",
			Expect: true,
		},
		{
			Name:   "nginx-otel",
			Input:  "Minimal Wolfi-based nginx HTTP server with OpenTelemetry observability integration for distributed tracing",
			Expect: true,
		},
		{
			Name:   "secretgen-controller",
			Input:  "secretgen-controller provides CRDs to specify what secrets need to be on Kubernetes cluster (to be generated or not)",
			Expect: true,
		},
		{
			Name:   "kubernetes-csi-external-snapshotter-fips",
			Input:  "Sidecar container that watches Kubernetes Snapshot CRD objects and triggers CreateSnapshot/DeleteSnapshot against a CSI endpoint. This container image includes the snapshotter, snapshot-controller and snapshot-validation-webhook.",
			Expect: true,
		},
		{
			Name:   "apache-tika-fips",
			Input:  "Apache Tika extracts metadata, text, and language from documents, enabling content analysis and indexing.",
			Expect: true,
		},
		{
			Name:   "request-4083",
			Input:  "nginx-stable image configured with Ruby for optimized application support",
			Expect: true,
		},
		{
			Name:   "mongodb-kubernetes-operator-fips",
			Input:  "Chainguard's MongoDB Kubernetes Operator image enables you to deploy a MongoDB community instance to a Kubernetes cluster, as well as support replica sets, scaling the replicas up or down, version upgrades, custom roles, and TLS security.",
			Expect: true,
		},
		{
			Name:   "clickhouse-iamguarded-fips",
			Input:  "FIPS-enabled minimal Wolfi-based ClickHouse analytics database image with IAMGuarded integration. [Clickhouse](https://clickhouse.com) is the fastest and most resource efficient open-source database for real-time apps and analytics.",
			Expect: true,
		},
		{
			Name:   "boky-postfix",
			Input:  "boky-postfix is an open-source Mail Transfer Agent that reliably sends and receives email. With Rspamd integration, it uses Rspamd’s built-in DKIM signing module to cryptographically sign outgoing emails, simplifying mail authentication and improving deliverability",
			Expect: true,
		},
		{
			Name:   "crossplane-function-auto-ready-fips",
			Input:  "This composition function automatically detects when composed resources are ready in Crossplane.",
			Expect: true,
		},
		{
			Name:   "vault",
			Input:  "Container image for Vault, a cross-platform secrets manager and authentication tool.",
			Expect: true,
		},
		{
			Name:   "gptscript",
			Input:  "Minimal [gptscript](https://github.com/gptscript-ai/gptscript) container image.",
			Expect: true,
		},
		{
			Name:   "maven",
			Input:  "Minimal image with the Maven build system.",
			Expect: true,
		},
		{
			Name:   "tetragon-fips",
			Input:  "FIPS-Complaint Images for Tetragon. eBPF-based Security Observability and Runtime Enforcement",
			Expect: true,
		},
		{
			Name:   "apache-kvrocks",
			Input:  "[Apache Kvrocks](https://kvrocks.apache.org/) is a distributed key-value database. Apache Kvrocks uses RocksDB as its storage engine and is compatible with Redis protocol.",
			Expect: true,
		},
		{
			Name:   "nginx-prometheus-exporter",
			Input:  "The `nginx-prometheus-exporter` image is designed to scrape metrics from an NGINX instance and expose them to Prometheus in a secure and minimal environment. Below are detailed instructions for using the image in both Docker and Kubernetes environments.",
			Expect: true,
		},
		{
			Name:   "velero",
			Input:  "Backup and migrate Kubernetes applications and their persistent volumes",
			Expect: true,
		},
		{
			Name:   "valkey",
			Input:  "Minimalist Wolfi-based [Valkey](https://github.com/valkey-io/valkey) image.",
			Expect: true,
		},
		{
			Name:   "k8s-metadata-injection",
			Input:  "Kubernetes metadata injection for New Relic APM to make a linkage between APM and Infrastructure data.",
			Expect: true,
		},
		{
			Name:   "gotenberg",
			Input:  "A developer-friendly API for converting numerous document formats into PDF files, and more!",
			Expect: true,
		},
		{
			Name:   "helm-chartmuseum",
			Input:  "Minimal image with [chartmuseum](https://github.com/helm/chartmuseum) binary.",
			Expect: true,
		},
		{
			Name:   "arangodb-fips",
			Input:  "ArangoDB is a native multi-model database with flexible data models for documents, graphs, and key-values. Build high performance applications using a convenient SQL-like query language or JavaScript extensions.",
			Expect: true,
		},
		{
			Name:   "dynamic-localpv-provisioner-fips",
			Input:  "Minimal Fips image of Dynamic LocalPV Provisioner an Kubernetes component that automates the provisioning of local persistent volumes.",
			Expect: true,
		},
		{
			Name:   "openai",
			Input:  "Minimal image with the OpenAI CLI.",
			Expect: true,
		},
		{
			Name:   "prometheus-operator",
			Input:  "Minimalist Wolfi-based image for Prometheus Operator. Prometheus Operator creates/configures/manages Prometheus clusters atop Kubernetes",
			Expect: true,
		},
		{
			Name:   "rsyslog",
			Input:  "[rsyslog](https://github.com/rsyslog/rsyslog) is a software utility used for log processing.",
			Expect: true,
		},
		{
			Name:   "cloud-provider-gcp-cloud-controller-manager-fips",
			Input:  "Kubernetes cloud controller manager for Google Cloud Platform (GCP), managing cloud-specific resources and integrations.",
			Expect: true,
		},
		{
			Name:   "code-server-fips",
			Input:  "VS Code in the browser",
			Expect: true,
		},
		{
			Name:   "tflint-fips",
			Input:  "A Pluggable Terraform Linter",
			Expect: true,
		},
		{
			Name:   "kube-metrics-adapter",
			Input:  "Minimal Adapter to expose custom metrics to Kubernetes HPA via Prometheus",
			Expect: true,
		},
		{
			Name:   "flux-operator-fips",
			Input:  "Flux Operator is a Kubernetes controller for managing the lifecycle of Flux CD",
			Expect: true,
		},
		{
			Name:   "request-6943",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "rust",
			Input:  "Minimal Wolfi-based Rust image for building Rust applications.",
			Expect: true,
		},
		{
			Name:   "klipper-helm",
			Input:  "Helm integration job image for K3s/RKE2 with automated chart lifecycle management",
			Expect: true,
		},
		{
			Name:   "authentik",
			Input:  "[Authentik](https://goauthentik.io/) is an open-source Identity Provider that provides single sign-on with support for SAML, OAuth2/OIDC, LDAP, and RADIUS protocols.",
			Expect: true,
		},
		{
			Name:   "kubernetes-autoscaler-addon-resizer-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "gatus-fips",
			Input:  "Gatus is a dev-oriented health dashboard that gives you the ability to monitor your services using HTTP, ICMP, TCP and DNS queries",
			Expect: true,
		},
		{
			Name:   "crossplane-provider-gitlab",
			Input:  "This image contains the Crossplane GitLab provider, which allows you to manage GitLab resources using Crossplane.",
			Expect: true,
		},
		{
			Name:   "tflint",
			Input:  "A Pluggable Terraform Linter",
			Expect: true,
		},
		{
			Name:   "omni-fips",
			Input:  "Omni managing Kubernetes clusters across bare metal, VMs, and cloud",
			Expect: true,
		},
		{
			Name:   "datadog-agent",
			Input:  "Minimalist Wolfi-based Datadog Agent to collect events and metrics from hosts and send them to Datadog.",
			Expect: true,
		},
		{
			Name:   "opentelemetry-collector-fips",
			Input:  "Minimal FIPS image with [opentelemetry-collector](https://github.com/open-telemetry/opentelemetry-collector).",
			Expect: true,
		},
		{
			Name:   "gha-runner-scale-set-controller",
			Input:  "Kubernetes controller for GitHub Actions self-hosted runners",
			Expect: true,
		},
		{
			Name:   "prism",
			Input:  "Prism is a set of packages for API mocking and contract testing with OpenAPI v2 (formerly known as Swagger) and OpenAPI v3.x.",
			Expect: true,
		},
		{
			Name:   "mcp-grafana",
			Input:  "mcp-grafana is a Model Context Protocol (MCP) server for Grafana that enables AI assistants and automation tools to interact with your Grafana dashboards, incidents, alerts, and datasources through a standardized protocol",
			Expect: true,
		},
		{
			Name:   "temporal-server",
			Input:  "Minimal image for [Temporal](https://docs.temporal.io/), a durable execution platform that handles intermittent failures and retries failed operations",
			Expect: true,
		},
		{
			Name:   "fluentd-iamguarded",
			Input:  "[Fluentd](https://www.fluentd.org/): Unified Logging Layer (project under CNCF)",
			Expect: true,
		},
		{
			Name:   "aws-otel-collector",
			Input:  "Distribution of OpenTelemetry Collector for sending data from EKS clusters to AWS monitoring services like CloudWatch and X-Ray.",
			Expect: true,
		},
		{
			Name:   "cerbos-fips",
			Input:  "Cerbos is the open core, language-agnostic, scalable authorization solution that makes user permissions and authorization simple to implement and manage by writing context-aware access control policies for your application resources.",
			Expect: true,
		},
		{
			Name:   "squid-proxy",
			Input:  "Squid Proxy is an open-source, high-performance, and highly configurable caching and forwarding web proxy. It is widely used for speeding up web servers by caching web, DNS, and other computer network lookups for a group of people sharing network resources, and for aiding security by filtering traffic.",
			Expect: true,
		},
		{
			Name:   "tritonserver-no-backend",
			Input:  "The Triton Inference Server provides an optimized cloud and edge inferencing solution.",
			Expect: true,
		},
		{
			Name:   "cluster-proportional-autoscaler",
			Input:  "Minimal Kubernetes Cluster Proportional Autoscaler Container",
			Expect: true,
		},
		{
			Name:   "quic-go",
			Input:  "A production-ready QUIC implementation in pure Go",
			Expect: true,
		},
		{
			Name:   "cloud-provider-azure-controller-manager-fips",
			Input:  "FIPS Controller manager for Azure CLI",
			Expect: true,
		},
		{
			Name:   "loki-fips",
			Input:  "This image contains the `loki` application for log aggregation. `loki` can be used to stream, aggregate, and query logs from apps and infrastructure.",
			Expect: true,
		},
		{
			Name:   "request-5754",
			Input:  "CLI for Databricks",
			Expect: true,
		},
		{
			Name:   "arangodb",
			Input:  "ArangoDB is a native multi-model database with flexible data models for documents, graphs, and key-values. Build high performance applications using a convenient SQL-like query language or JavaScript extensions.",
			Expect: true,
		},
		{
			Name:   "ratify",
			Input:  "Artifact Ratification Framework (CNCF Sandbox)",
			Expect: true,
		},
		{
			Name:   "contour",
			Input:  "Contour is an ingress controller for Kubernetes that works by deploying the Envoy proxy as a reverse proxy and load balancer. Contour supports dynamic configuration updates out of the box while maintaining a lightweight profile.",
			Expect: true,
		},
		{
			Name:   "pgadmin4",
			Input:  "pgAdmin is an open source administration and development platform for PostgreSQL.",
			Expect: true,
		},
		{
			Name:   "step-issuer-fips",
			Input:  "Minimal FIPS image of [step-issuer](https://smallstep.com/docs/platform/), a certificate issuer for cert-manager using step-ca as a backend CA.",
			Expect: true,
		},
		{
			Name:   "perl",
			Input:  "Container image for building Perl applications.",
			Expect: true,
		},
		{
			Name:   "secrets-store-csi-driver-provider-gcp",
			Input:  "Minimal image with the Kubernetes Secrets Store CSI Driver GCP Plugin.",
			Expect: true,
		},
		{
			Name:   "chainguard-source",
			Input:  "Fetches all sources referenced by a Chainguard Package or Image SBOM.",
			Expect: true,
		},
		{
			Name:   "authservice",
			Input:  "Move OIDC token acquisition out of your app code and into the Istio mesh",
			Expect: true,
		},
		{
			Name:   "nginx-otel-fips",
			Input:  "Minimal Wolfi-based nginx HTTP server with OpenTelemetry observability integration and FIPS 140-2 hardened cryptography",
			Expect: true,
		},
		{
			Name:   "amazon-cloudwatch-agent",
			Input:  "CloudWatch Agent enables you to collect and export host-level metrics and logs on instances running Linux or Windows server. ",
			Expect: true,
		},
		{
			Name:   "kube-conformance",
			Input:  "Minimal container image for running Kubernetes conformance tests",
			Expect: true,
		},
		{
			Name:   "erlang",
			Input:  "Container image for building Erlang applications.",
			Expect: true,
		},
		{
			Name:   "amazon-k8s-cni",
			Input:  "Networking plugin repository for pod networking in Kubernetes using Elastic Network Interfaces on AWS",
			Expect: true,
		},
		{
			Name:   "kubernetes-csi-driver-nfs-fips",
			Input:  "This driver allows Kubernetes to access an NFS server running on a Linux node. This image is FIPS compliant.",
			Expect: true,
		},
		{
			Name:   "nova",
			Input:  "Nova is a cli tool to find outdated or deprecated Helm charts running in your Kubernetes cluster.",
			Expect: true,
		},
		{
			Name:   "jre-crac",
			Input:  "Minimalist Wolfi-based OpenJDK JRE image with [CRaC](https://openjdk.org/projects/crac/) support. Used for running Java applications.",
			Expect: true,
		},
		{
			Name:   "graphicsmagick",
			Input:  "Minimal container image with GraphicsMagick, a collection of tools allowing you to read, write, and manipulate images in a variety of formats.",
			Expect: true,
		},
		{
			Name:   "git-sync-fips",
			Input:  "A sidecar app which clones a git repo and keeps it in sync with the upstream.",
			Expect: true,
		},
		{
			Name:   "liberica-jre",
			Input:  "Free and open source Progressive Java Runtime for modern Java deployments",
			Expect: true,
		},
		{
			Name:   "percona-server-mongodb-operator",
			Input:  "Minimal Kubernetes Operator image which deploys and manages Percona Server for MongoDB on Kubernetes clusters.",
			Expect: true,
		},
		{
			Name:   "regclient",
			Input:  "regclient is a client interface to OCI conformant registries and content shipped with the OCI Image Layout",
			Expect: true,
		},
		{
			Name:   "cert-manager",
			Input:  "Wolfi-based images for [cert-manager](https://cert-manager.io), a tool for provisioning and managing TLS certificates in Kubernetes.",
			Expect: true,
		},
		{
			Name:   "eks-distro",
			Input:  "An open-source distribution of Kubernetes from AWS",
			Expect: true,
		},
		{
			Name:   "cloud-sql-proxy-fips",
			Input:  "The Cloud SQL Auth Proxy is a utility for ensuring secure connections to Cloud SQL instances.",
			Expect: true,
		},
		{
			Name:   "falcoctl",
			Input:  "Minimalist Wolfi-based image for `falcoctl`.",
			Expect: true,
		},
		{
			Name:   "distribution",
			Input:  "The toolkit to pack, ship, store, and deliver container content",
			Expect: true,
		},
		{
			Name:   "prometheus-iamguarded",
			Input:  "Prometheus is a monitoring system and time series database",
			Expect: true,
		},
		{
			Name:   "frr-fips",
			Input:  "The FRRouting Protocol Suite",
			Expect: true,
		},
		{
			Name:   "octo-sts",
			Input:  "A GitHub App that acts like a Security Token Service (STS) for the Github API.",
			Expect: true,
		},
		{
			Name:   "cloudflared-fips",
			Input:  "Cloudflare Tunnel client (formerly Argo Tunnel)",
			Expect: true,
		},
		{
			Name:   "kubeflow-pipelines-visualization-server",
			Input:  "Minimal image with [ml-pipeline/visualization-server](https://github.com/kubeflow/pipelines/tree/master/backend/src/apiserver/visualization).",
			Expect: true,
		},
		{
			Name:   "kong",
			Input:  "Kong is a Cloud-Native API Gateway and AI Gateway",
			Expect: true,
		},
		{
			Name:   "dask-kubernetes",
			Input:  "Native Kubernetes integration for Dask",
			Expect: true,
		},
		{
			Name:   "vault-secrets-operator",
			Input:  "Vault Secrets Operator (VSO) allows Pods to consume Vault secrets natively from Kubernetes Secrets.",
			Expect: true,
		},
		{
			Name:   "headlamp-fips",
			Input:  "Minimal Fips image of Headlamp is an easy-to-use and extensible Kubernetes web UI designed for developers and cluster operators.",
			Expect: true,
		},
		{
			Name:   "cerbos",
			Input:  "Cerbos is the open core, language-agnostic, scalable authorization solution that makes user permissions and authorization simple to implement and manage by writing context-aware access control policies for your application resources.",
			Expect: true,
		},
		{
			Name:   "valkey-iamguarded",
			Input:  "Valkey is an open source, in-memory data store used by millions of developers as a cache, vector database, document database, streaming engine, and message broker.",
			Expect: true,
		},
		{
			Name:   "crossplane-function-auto-ready",
			Input:  "This composition function automatically detects when composed resources are ready in Crossplane.",
			Expect: true,
		},
		{
			Name:   "pluto",
			Input:  "A cli tool to help discover deprecated apiVersions in Kubernetes",
			Expect: true,
		},
		{
			Name:   "knative-operator",
			Input:  "Combined operator for Knative.",
			Expect: true,
		},
		{
			Name:   "rekor",
			Input:  "Rekor is one of the core components of the sigstore stack.",
			Expect: true,
		},
		{
			Name:   "sriov-network-device-plugin",
			Input:  "SRIOV network device plugin for Kubernetes",
			Expect: true,
		},
		{
			Name:   "gatekeeper",
			Input:  "Minimal [Gatekeeper](https://open-policy-agent.github.io/gatekeeper) image for enforcing Kubernetes policies using Open Policy Agent",
			Expect: true,
		},
		{
			Name:   "spicedb",
			Input:  "[SpiceDB](https://authzed.com/spicedb) is an open-source authorization database inspired by Google's Zanzibar, providing scalable and fine-grained access control for applications.",
			Expect: true,
		},
		{
			Name:   "actions-runner",
			Input:  "actions-runner is a self-hosted application that runs GitHub Actions jobs on your own infrastructure.",
			Expect: true,
		},
		{
			Name:   "hydra",
			Input:  "Ory Hydra is a hardened, OpenID Certified OAuth 2.0 Server and OpenID Connect Provider optimized for low-latency, high throughput, and low resource consumption.",
			Expect: true,
		},
		{
			Name:   "polaris-fips",
			Input:  "FIPS-compliant version of Polaris. Polaris is an open source policy engine for Kubernetes that validates and remediates resource configuration.",
			Expect: true,
		},
		{
			Name:   "aws-s3-controller",
			Input:  "S3 controller is an ACK service controller for Amazon Simple Storage Service (S3).",
			Expect: true,
		},
		{
			Name:   "metrics-agent-fips",
			Input:  "metrics-agent collects Kubernetes allocation and utilization data",
			Expect: true,
		},
		{
			Name:   "crossplane-provider-terraform",
			Input:  "Crossplane Terraform Provider enables provisioning and managing infrastructure using Terraform within a Crossplane control plane.",
			Expect: true,
		},
		{
			Name:   "teleport",
			Input:  "Teleport is an access management platform designed to provide secure and unified access to various infrastructure resources such as SSH, Kubernetes clusters, databases, and web applications",
			Expect: true,
		},
		{
			Name:   "cluster-api-aws-controller-fips",
			Input:  "Kubernetes Cluster API Provider AWS provides consistent deployment and day 2 operations of \"self-managed\" and EKS Kubernetes clusters on AWS",
			Expect: true,
		},
		{
			Name:   "kuberay-operator-fips",
			Input:  "A toolkit to run Ray applications on Kubernetes",
			Expect: true,
		},
		{
			Name:   "request-4952",
			Input:  "Container image for running the AWX Kubernetes operator",
			Expect: true,
		},
		{
			Name:   "crossplane-gcp",
			Input:  "Crossplane GCP Providers deliver Kubernetes-native APIs for provisioning and managing Google Cloud resources through Crossplane.",
			Expect: true,
		},
		{
			Name:   "k6-operator-fips",
			Input:  "FIPS-enabled Kubernetes operator for running distributed k6 performance tests",
			Expect: true,
		},
		{
			Name:   "mariadb-operator-fips",
			Input:  "MariaDB operator FIPS image is a FIPS-compliant image of the MariaDB operator, It automates the deployment, scaling, and management of MariaDB instances in Kubernetes clusters, providing declarative configuration and lifecycle management.",
			Expect: true,
		},
		{
			Name:   "yara",
			Input:  "The pattern matching swiss knife.",
			Expect: true,
		},
		{
			Name:   "kibana-iamguarded",
			Input:  "Your window into the Elastic Stack.",
			Expect: true,
		},
		{
			Name:   "cypress-base",
			Input:  "Minimal image for cypress/base.",
			Expect: true,
		},
		{
			Name:   "pgbouncer-iamguarded",
			Input:  "This image contains the CLI for the [pgbouncer](https://www.pgbouncer.org/) connection pooler for PostgreSQL with IAMGuarded support.",
			Expect: true,
		},
		{
			Name:   "pgpool2",
			Input:  "Middleware that works between PostgreSQL servers and a PostgreSQL database client.",
			Expect: true,
		},
		{
			Name:   "promtail",
			Input:  "This image contains the `promtail` application for log aggregation. `promtail` is the log aggregator that ships logs to Loki and/or Prometheus. It runs as an agent and scrapes logs from files, containers, and hosts and ships them to a logging backend.",
			Expect: true,
		},
		{
			Name:   "rook-ceph",
			Input:  "Storage Orchestration for Kubernetes. This is specifically for the Rook Ceph operator, which provides storage solutions for Kubernetes clusters.",
			Expect: true,
		},
		{
			Name:   "rke2-cloud-provider",
			Input:  "rke2-cloud-provider image.",
			Expect: true,
		},
		{
			Name:   "graalvm",
			Input:  "GraalVM is an advanced JDK with ahead-of-time Native Image compilation",
			Expect: true,
		},
		{
			Name:   "openssh-server-fips",
			Input:  "OpenSSH Server FIPS is a secure shell (SSH) server implementation that provides encrypted communication between clients and servers, with FIPS compliance.",
			Expect: true,
		},
		{
			Name:   "secrets-store-csi-driver-provider-azure-fips",
			Input:  "Azure Key Vault provider for Secret Store CSI driver allows you to get secret contents stored in Azure Key Vault instance and use the Secret Store CSI driver interface to mount them into Kubernetes pods.",
			Expect: true,
		},
		{
			Name:   "tritonserver-vllm-backend-fips",
			Input:  "The Triton Inference Server provides an optimized cloud and edge inferencing solution with vllm backend",
			Expect: true,
		},
		{
			Name:   "aws-gateway-controller-fips",
			Input:  "AWS Application Networking is an implementation of the Kubernetes Gateway API.",
			Expect: true,
		},
		{
			Name:   "ipfs-cluster-fips",
			Input:  "Pinset orchestration for IPFS",
			Expect: true,
		},
		{
			Name:   "request-5675",
			Input:  "Lightweight DICOM server for medical imaging with RESTful API and web interface",
			Expect: true,
		},
		{
			Name:   "cert-manager-startupapicheck",
			Input:  "Automatically provision and manage TLS certificates in Kubernetes",
			Expect: true,
		},
		{
			Name:   "newrelic-prometheus-configurator",
			Input:  "Minimal [newrelic-prometheus-configurator](https://github.com/newrelic/newrelic-prometheus-configurator) container image.",
			Expect: true,
		},
		{
			Name:   "kubernetes-replicator-fips",
			Input:  "[kubernetes-replicator](https://github.com/mittwald/kubernetes-replicator) is a custom Kubernetes controller that can be used to make secrets and config maps available in multiple namespaces.",
			Expect: true,
		},
		{
			Name:   "velero-plugin-for-microsoft-azure",
			Input:  "Velero plugin for Microsoft Azure that provides backup and restore functionality for Azure Blob Storage and Azure Disk snapshots",
			Expect: true,
		},
		{
			Name:   "opentelemetry-nodejs-instrumentation",
			Input:  "Auto instrumention for any Node application to capture telemetry from a number of popular libraries and frameworks.",
			Expect: true,
		},
		{
			Name:   "cluster-api-helm-controller",
			Input:  "CAAPH uses Helm charts to manage the installation and lifecycle of Cluster API add-ons.",
			Expect: true,
		},
		{
			Name:   "mlflow",
			Input:  "A minimal, [Wolfi](https://github.com/wolfi-dev)-based image for MLflow, an open source platform for the machine learning lifecycle.",
			Expect: true,
		},
		{
			Name:   "kyverno-reports-server",
			Input:  "Reports server provides a scalable solution for storing policy reports and cluster policy reports. It moves reports out of etcd and stores them in a PostgreSQL database instance.",
			Expect: true,
		},
		{
			Name:   "aws-node-termination-handler",
			Input:  "Gracefully handle EC2 instance shutdown within Kubernetes",
			Expect: true,
		},
		{
			Name:   "postgres-repmgr-iamguarded",
			Input:  "PostgreSQL HA, a cluster solution using the PostgreSQL replication manager.",
			Expect: true,
		},
		{
			Name:   "grafana-agent-operator",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "elixir-fips",
			Input:  "Container image for building Elixir applications with FIPS",
			Expect: true,
		},
		{
			Name:   "keycloak-operator-fips",
			Input:  "A Kubernetes Operator based on the Operator SDK for installing and managing Keycloak with FIPS support.",
			Expect: true,
		},
		{
			Name:   "zipkin-slim",
			Input:  "Smaller distribution of Zipkin which supports Elasticsearch storage and HTTP or gRPC span collection",
			Expect: true,
		},
		{
			Name:   "rstudio",
			Input:  "Minimal [RStudio](https://github.com/rstudio/rstudio) container image.",
			Expect: true,
		},
		{
			Name:   "go-msft-fips",
			Input:  "Container image for building Go applications with FIPS",
			Expect: true,
		},
		{
			Name:   "nextflow",
			Input:  "Nextflow is a domain-specific language (DSL) for data-driven computational pipelines.",
			Expect: true,
		},
		{
			Name:   "rsyslog-fips",
			Input:  "[rsyslog](https://github.com/rsyslog/rsyslog) is a software utility used for log processing.",
			Expect: true,
		},
		{
			Name:   "komodo",
			Input:  "Chainguard images for [Komodo](https://github.com/moghtech/komodo), a DevOps platform for building and deploying software across multiple servers.",
			Expect: true,
		},
		{
			Name:   "spark-fips",
			Input:  "Spark provides high-level APIs in Scala, Java, Python, and R, and an optimized engine that supports general computation graphs for data analysis.",
			Expect: true,
		},
		{
			Name:   "glibc-dynamic",
			Input:  "Base image with just enough to run arbitrary glibc binaries.",
			Expect: true,
		},
		{
			Name:   "kserve",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "falcosidekick",
			Input:  "Minimalist Wolfi-based image for `falcosidekick`.",
			Expect: true,
		},
		{
			Name:   "sigstore-policy-controller-fips",
			Input:  "Fips version of Policy Controller image that is part of the Sigstore stack",
			Expect: true,
		},
		{
			Name:   "kibana-iamguarded-fips",
			Input:  "Your window into the Elastic Stack.",
			Expect: true,
		},
		{
			Name:   "prometheus-node-exporter-iamguarded",
			Input:  "Prometheus exporter for hardware and OS metrics exposed by LINUX kernels",
			Expect: true,
		},
		{
			Name:   "sealed-secrets-controller-fips",
			Input:  "A Kubernetes controller and tool for one-way encrypted Secrets, enabling safe GitOps-friendly secret management.",
			Expect: true,
		},
		{
			Name:   "gha-runner-scale-set-controller-fips",
			Input:  "Kubernetes controller for GitHub Actions self-hosted runners",
			Expect: true,
		},
		{
			Name:   "kubernetes-csi-external-attacher",
			Input:  "Chainguard image for [kubernetes-csi/external-attacher](https://github.com/kubernetes-csi/external-attacher). Watches Kubernetes VolumeAttachment objects and triggers ControllerPublish/Unpublish against a CSI endpoint.",
			Expect: true,
		},
		{
			Name:   "opentelemetry-python-instrumentation",
			Input:  "OpenTelemetry auto-instrumentation and instrumentation libraries for Python",
			Expect: true,
		},
		{
			Name:   "pypiserver",
			Input:  "Minimal PyPI server for uploading & downloading packages with pip/easy_install",
			Expect: true,
		},
		{
			Name:   "image-factory",
			Input:  "Image Factory is a flexible artifact-build service for Talos Linux, offering automated creation of ISOs, disks, cloud images, and installers and ensuring reproductible builds.",
			Expect: true,
		},
		{
			Name:   "k8ssandra-operator",
			Input:  "The Kubernetes operator for K8ssandra",
			Expect: true,
		},
		{
			Name:   "rabbitmq",
			Input:  "[RabbitMQ](https://github.com/rabbitmq/rabbitmq-server) is a message broker.",
			Expect: true,
		},
		{
			Name:   "metrics-server-iamguarded",
			Input:  "Metrics Server is a Kubernetes component that collects and provides resource usage metrics (CPU, memory) for nodes and pods",
			Expect: true,
		},
		{
			Name:   "grafana",
			Input:  "A minimal wolfi-based image for grafana, which is an open-source monitoring and observability application",
			Expect: true,
		},
		{
			Name:   "opensearch-k8s-operator",
			Input:  "The Kubernetes OpenSearch Operator is used for automating the deployment, provisioning, management, and orchestration of OpenSearch clusters and OpenSearch dashboards.",
			Expect: true,
		},
		{
			Name:   "prometheus-alertmanager",
			Input:  "Minimalist Wolfi-based image for Prometheus Alertmanager. Handles alerts sent by client applications such as the Prometheus server. It takes care of deduplicating, grouping, and routing to the correct receiver.",
			Expect: true,
		},
		{
			Name:   "mailpit",
			Input:  "Mailpit is an email and SMTP testing tool with API for developers.",
			Expect: true,
		},
		{
			Name:   "kube-logging-operator-node-exporter",
			Input:  "Custom runner based Prometheus node exporter for kube-logging logging-operator",
			Expect: true,
		},
		{
			Name:   "json-mock",
			Input:  "Container image for json-server to mock REST/JSON APIs",
			Expect: true,
		},
		{
			Name:   "crossplane-function-go-templating-fips",
			Input:  "This composition function allows you to compose Crossplane resources using Go templates.",
			Expect: true,
		},
		{
			Name:   "laravel",
			Input:  "Minimalist Wolfi-based Laravel images for developing, building, and running Laravel applications.",
			Expect: true,
		},
		{
			Name:   "mattermost-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "aws-otel-collector-fips",
			Input:  "Distribution of OpenTelemetry Collector for sending data from EKS clusters to AWS monitoring services like CloudWatch and X-Ray.",
			Expect: true,
		},
		{
			Name:   "rabbitmq-messaging-topology-operator-iamguarded-fips",
			Input:  "RabbitMQ messaging topology operator",
			Expect: true,
		},
		{
			Name:   "smarter-device-manager",
			Input:  "Minimalist Wolfi-based image for smarter device manager.",
			Expect: true,
		},
		{
			Name:   "custom-pod-autoscaler-operator",
			Input:  "Operator for managing Kubernetes Custom Pod Autoscalers (CPA)",
			Expect: true,
		},
		{
			Name:   "prometheus-yet-another-cloudwatch-exporter-fips",
			Input:  "A FIPS-compliant Prometheus exporter for AWS CloudWatch.",
			Expect: true,
		},
		{
			Name:   "fluent-operator",
			Input:  "Operator for Fluent Bit and Fluentd - previously known as FluentBit Operator",
			Expect: true,
		},
		{
			Name:   "bun",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "xeol",
			Input:  "A scanner for end-of-life (EOL) software and dependencies in container images, filesystems, and SBOMs",
			Expect: true,
		},
		{
			Name:   "spark",
			Input:  "Spark provides high-level APIs in Scala, Java, Python, and R, and an optimized engine that supports general computation graphs for data analysis.",
			Expect: true,
		},
		{
			Name:   "gitlab-fips",
			Input:  "FIPS-compliant GitLab images providing a complete DevOps platform that meets Federal Information Processing Standards for cryptographic operations, source code management, CI/CD automation, and collaboration tools.",
			Expect: true,
		},
		{
			Name:   "minio",
			Input:  "Minimal image with Minio.",
			Expect: true,
		},
		{
			Name:   "orthanc",
			Input:  "Lightweight DICOM server for medical imaging with RESTful API and web interface",
			Expect: true,
		},
		{
			Name:   "request-7571",
			Input:  "custom image with monit package",
			Expect: true,
		},
		{
			Name:   "images/povray-fips",
			Input:  "POV-Ray is a ray-tracing program that generates images from text-based scene descriptions, with FIPS 140-3 compliant cryptographic modules.",
			Expect: true,
		},
		{
			Name:   "cert-exporter-fips",
			Input:  "A minimal FIPS image for cert-exporter: an application that exports certificate expiration metrics from disk, Kubernetes, and AWS Secrets Manager to Prometheus.",
			Expect: true,
		},
		{
			Name:   "multus-cni-fips",
			Input:  "A CNI meta-plugin for multi-homed pods in Kubernetes",
			Expect: true,
		},
		{
			Name:   "flux-helm-controller-iamguarded",
			Input:  "minimal zero CVE flux images",
			Expect: true,
		},
		{
			Name:   "google-cloud-sdk",
			Input:  "Minimal image with the [Google Cloud SDK](https://cloud.google.com/sdk/).",
			Expect: true,
		},
		{
			Name:   "apm-server",
			Input:  "Elastic APM is an application performance monitoring system built on the Elastic Stack.",
			Expect: true,
		},
		{
			Name:   "airflow",
			Input:  "A minimal, [wolfi](https://github.com/wolfi-dev)-based image for Apache Airflow.[Apache Airflow](https://github.com/apache/airflow) is a platform to programmatically author, schedule, and monitor workflows.",
			Expect: true,
		},
		{
			Name:   "kube-logging-operator",
			Input:  "Minimal Logging operator for Kubernetes Image",
			Expect: true,
		},
		{
			Name:   "skopeo-fips",
			Input:  "Minimalist Wolfi-based skopeo image for interacting with container registries.",
			Expect: true,
		},
		{
			Name:   "rabbitmq-cluster-operator-iamguarded-fips",
			Input:  "RabbitMQ Cluster Kubernetes Operator",
			Expect: true,
		},
		{
			Name:   "local-volume-provisioner-fips",
			Input:  "Static provisioner of local volumes",
			Expect: true,
		},
		{
			Name:   "mongodb-iamguarded",
			Input:  "[MongoDB](https://www.mongodb.com/) is a document-oriented database management system. MongoDB is a popular example of a NoSQL database, and stores data in JSON-like documents.",
			Expect: true,
		},
		{
			Name:   "pgwatch",
			Input:  "pgwatch image.",
			Expect: true,
		},
		{
			Name:   "kube-rbac-proxy",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "sealed-secrets-controller",
			Input:  "A Kubernetes controller and tool for one-way encrypted Secrets, enabling safe GitOps-friendly secret management.",
			Expect: true,
		},
		{
			Name:   "php",
			Input:  "Minimalist Wolfi-based PHP images for building and running PHP applications. Includes both `dev` and `fpm` variants.",
			Expect: true,
		},
		{
			Name:   "pytorch-fips",
			Input:  "A minimal, [wolfi](https://github.com/wolfi-dev)-based FIPS compliant image for pytorch, a Python package that provides two high-level features: Tensor computation with strong GPU acceleration and Deep neural networks built on a tape-based autograd system.",
			Expect: true,
		},
		{
			Name:   "cassandra-fips",
			Input:  "[Cassandra](https://cassandra.apache.org) is a free and open-source, distributed, wide-column store, NoSQL database.",
			Expect: true,
		},
		{
			Name:   "prometheus-elasticsearch-exporter-iamguarded",
			Input:  "Minimalist Wolfi-based Prometheus Elasticsearch Exporter image for exporting various metrics about Elasticsearch.",
			Expect: true,
		},
		{
			Name:   "kong-fips",
			Input:  "Kong is a Cloud-Native API Gateway and AI Gateway",
			Expect: true,
		},
		{
			Name:   "prometheus-statsd-exporter-iamguarded",
			Input:  "[statsd_exporter](https://github.com/prometheus/statsd_exporter) receives StatsD-style metrics and exports them as Prometheus metrics.",
			Expect: true,
		},
		{
			Name:   "flannel",
			Input:  "Flannel is a network fabric for Kubernetes, providing a way to manage network configurations across a cluster.",
			Expect: true,
		},
		{
			Name:   "fulcio",
			Input:  "Minimal Fulcio image",
			Expect: true,
		},
		{
			Name:   "azure-ipam",
			Input:  "Azure VNET IPAM plugins manage IP address assignments to containers.",
			Expect: true,
		},
		{
			Name:   "spire",
			Input:  "Minimalist Wolfi-based `spire` images.",
			Expect: true,
		},
		{
			Name:   "azure-functions-python",
			Input:  "Azure Functions is a managed platform-as-a-service provider which offers scalable and serverless hosting for Python code projects. It extends the Azure platform with the capability to implement code triggered by many events occurring in Azure, on-premises or other 3rd party services.",
			Expect: true,
		},
		{
			Name:   "envoy",
			Input:  "[Envoy](https://www.envoyproxy.io/) Cloud-native high-performance edge/middle/service proxy",
			Expect: true,
		},
		{
			Name:   "fluent-bit",
			Input:  "[Fluent Bit](https://fluentbit.io) is a lightweight and high performance log processor.",
			Expect: true,
		},
		{
			Name:   "sigstore-scaffolding",
			Input:  "Minimal Wolfi-based [Sigstore](https://sigstore.dev) images.",
			Expect: true,
		},
		{
			Name:   "crossplane-keycloak-fips",
			Input:  "A Crossplane provider for Keycloak FIPS version",
			Expect: true,
		},
		{
			Name:   "TODO Find Out The Name",
			Input:  "Minimal Wolfi-based Java JDK image using Adoptium OpenJDK. Used for compiling Java applications.",
			Expect: true,
		},
		{
			Name:   "elasticsearch-fips",
			Input:  "Elasticsearch is a distributed search and analytics engine, scalable data store and vector database optimized for speed and relevance on production-scale workloads.",
			Expect: true,
		},
		{
			Name:   "kapp-controller-fips",
			Input:  "Continuous delivery and package management for Kubernetes.",
			Expect: true,
		},
		{
			Name:   "keda",
			Input:  "Minimal image with the Keda binary.",
			Expect: true,
		},
		{
			Name:   "newrelic-k8s-events-forwarder-fips",
			Input:  "A FIPS-Compliant & lightweight Kubernetes event forwarder that streams cluster events to New Relic for centralized monitoring and analysis.",
			Expect: true,
		},
		{
			Name:   "prometheus-postgres-exporter-iamguarded",
			Input:  "A PostgreSQL metric exporter for Prometheus",
			Expect: true,
		},
		{
			Name:   "gitea-fips",
			Input:  "A painless self-hosted all-in-one software development service, including Git hosting, code review, team collaboration, package registry and CI/CD",
			Expect: true,
		},
		{
			Name:   "kustomize-mutating-webhook",
			Input:  "A dynamic solution to patch FluxCD Kustomization resources, seamlessly integrating and federating substitution variables across multiple namespaces.",
			Expect: true,
		},
		{
			Name:   "longhorn",
			Input:  "A lightweight, reliable distributed block storage system for Kubernetes.",
			Expect: true,
		},
		{
			Name:   "zookeeper-iamguarded",
			Input:  "[Apache ZooKeeper](https://zookeeper.apache.org/) is an effort to develop and maintain an open-source server which enables highly reliable distributed coordination.",
			Expect: true,
		},
		{
			Name:   "haproxy-ingress",
			Input:  "Kubernetes ingress controller implementation for HAProxy",
			Expect: true,
		},
		{
			Name:   "vertical-pod-autoscaler",
			Input:  "Image to automatically adjust the amount of CPU and memory requested by pods running in the Kubernetes Cluster",
			Expect: true,
		},
		{
			Name:   "amazon-cloudwatch-agent-operator",
			Input:  "The Amazon CloudWatch Agent Operator is software developed to manage the CloudWatch Agent on kubernetes.",
			Expect: true,
		},
		{
			Name:   "hubble-export-stdout-fips",
			Input:  "FIPS-compliant hubble-export-stdout exports Hubble data to stdout.",
			Expect: true,
		},
		{
			Name:   "podinfo",
			Input:  "Podinfo is a tiny web application that provides Go microsrvice template for kubernetes",
			Expect: true,
		},
		{
			Name:   "livekit-server-fips",
			Input:  "livekit-server is an open-source media server for real-time audio, video, and data, designed for low latency and scalability",
			Expect: true,
		},
		{
			Name:   "kiali-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "request-2682-fips",
			Input:  "A FIPS-complaint container registry offering secure image storage, access control, scanning, and replication, built as Harbor's customized Docker Distribution registry service.",
			Expect: true,
		},
		{
			Name:   "kyverno-policy-reporter-fips",
			Input:  "Monitoring and Observability Tool for the [PolicyReport CRD](https://kyverno.github.io/policy-reporter/) with an optional UI.",
			Expect: true,
		},
		{
			Name:   "metrics-agent",
			Input:  "metrics-agent collects Kubernetes allocation and utilization data",
			Expect: true,
		},
		{
			Name:   "cloud-provider-vsphere",
			Input:  "Minimal image of the [Kubernetes cloud provider interface for vSphere](https://cloud-provider-vsphere.sigs.k8s.io/)",
			Expect: true,
		},
		{
			Name:   "spire-controller-manager-fips",
			Input:  "The SPIRE Controller Manager provides automated workload identity management for Kubernetes clusters through SPIRE",
			Expect: true,
		},
		{
			Name:   "logstash-oss-with-opensearch-output-plugin",
			Input:  "An image with the Logstash plugin that sends event data to a OpenSearch clusters and stores as an index.",
			Expect: true,
		},
		{
			Name:   "solr",
			Input:  "Solr is an open-source multi-modal search platform built on top of Lucene.",
			Expect: true,
		},
		{
			Name:   "spark-iamguarded",
			Input:  "Apache Spark, is a multi-language engine for executing data engineering, data science, and machine learning on single-node machines or clusters.",
			Expect: true,
		},
		{
			Name:   "prometheus-pushgateway",
			Input:  "Minimal Prometheus Pushgateway Image",
			Expect: true,
		},
		{
			Name:   "cluster-api-vsphere-controller",
			Input:  "Kubernetes-native declarative infrastructure for vSphere",
			Expect: true,
		},
		{
			Name:   "azure-workload-identity-webhook",
			Input:  "A webhook for Kubernetes that enables Azure Active Directory based authentication from Kubernetes workloads to Azure resources.",
			Expect: true,
		},
		{
			Name:   "superset-iamguarded",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "apache-exporter",
			Input:  "apache-exporter exposes Apache HTTP server metrics from mod_status in Prometheus format for monitoring",
			Expect: true,
		},
		{
			Name:   "dynamic-localpv-provisioner",
			Input:  "Dynamic Local Volumes for Kubernetes Stateful workloads.",
			Expect: true,
		},
		{
			Name:   "gcp-compute-persistent-disk-csi-driver-fips",
			Input:  "The Google Compute Engine Persistent Disk (GCE PD) Container Storage Interface (CSI) Storage Plugin.",
			Expect: true,
		},
		{
			Name:   "pg-timetable-fips",
			Input:  "An advanced standalone job scheduler for PostgreSQL, offering many advantages over traditional schedulers such as cron and others. This image provides FIPS support for pg-timetable.",
			Expect: true,
		},
		{
			Name:   "seaweedfs-operator-fips",
			Input:  "seaweedfs kubernetes operator",
			Expect: true,
		},
		{
			Name:   "request-5754",
			Input:  "CLI for Databricks",
			Expect: true,
		},
		{
			Name:   "azure-functions-python-fips",
			Input:  "Azure Functions is a managed platform-as-a-service provider which offers scalable and serverless hosting for Python code projects. It extends the Azure platform with the capability to implement code triggered by many events occurring in Azure, on-premises or other 3rd party services.",
			Expect: true,
		},
		{
			Name:   "percona-server-mongodb-operator-fips",
			Input:  "Minimal FIPS compliant Kubernetes Operator image which deploys and manages Percona Server for MongoDB on Kubernetes clusters.",
			Expect: true,
		},
		{
			Name:   "prometheus-node-exporter-iamguarded-fips",
			Input:  "Prometheus exporter for hardware and OS metrics exposed by LINUX kernels",
			Expect: true,
		},
		{
			Name:   "gradle",
			Input:  "Chainguard Image with [Gradle](https://gradle.org/), an open source build system for Java, Android, and Kotlin.",
			Expect: true,
		},
		{
			Name:   "bash",
			Input:  "Container image with only Bash and libc. Suitable for running any small scripts or binaries that need Bash instead of the BusyBox shell.",
			Expect: true,
		},
		{
			Name:   "objectstorage-fips",
			Input:  "Container Object Storage Interface (COSI) controller and sidecar",
			Expect: true,
		},
		{
			Name:   "rqlite-fips",
			Input:  "Minimal image with rqlite, a lightweight, distributed relational database built on SQLite. rqlite uses the Raft consensus protocol to provide strong consistency across a cluster of nodes, making it an ideal solution for lightweight, fault-tolerant distributed databases.",
			Expect: true,
		},
		{
			Name:   "apache-beam-python-sdk",
			Input:  "Apache Beam is a unified programming model for Batch and Streaming data processing.",
			Expect: true,
		},
		{
			Name:   "kaniko",
			Input:  "Build Container Images In Kubernetes",
			Expect: true,
		},
		{
			Name:   "kaniko-warmer",
			Input:  "Build Container Images In Kubernetes",
			Expect: true,
		},
		{
			Name:   "cinc-auditor",
			Input:  "Open source toolkit for applying Chef Inspec audit and test profiles.",
			Expect: true,
		},
		{
			Name:   "dotnet",
			Input:  "Minimal container image for .NET and the .NET Tools.",
			Expect: true,
		},
		{
			Name:   "aws-for-fluent-bit",
			Input:  "Minimal [aws-for-fluent-bit](https://github.com/aws/aws-for-fluent-bit) Image",
			Expect: true,
		},
		{
			Name:   "request-5891",
			Input:  "Highly Available PostgreSQL cluster using Docker",
			Expect: true,
		},
		{
			Name:   "sonar-scanner-cli",
			Input:  "Scanner CLI for SonarQube and SonarCloud",
			Expect: true,
		},
		{
			Name:   "nats-box",
			Input:  "A lightweight container with NATS utilities.",
			Expect: true,
		},
		{
			Name:   "influxdb",
			Input:  "Minimal image with influxdb.",
			Expect: true,
		},
		{
			Name:   "deck",
			Input:  "deck is a command-line interface for managing Kong Gateway configurations declaratively",
			Expect: true,
		},
		{
			Name:   "cloud-sql-proxy",
			Input:  "The Cloud SQL Auth Proxy is a utility for ensuring secure connections to Cloud SQL instances.",
			Expect: true,
		},
		{
			Name:   "mongodb",
			Input:  "[MongoDB](https://www.mongodb.com/) is a document-oriented database management system. MongoDB is a popular example of a NoSQL database, and stores data in JSON-like documents.",
			Expect: true,
		},
		{
			Name:   "rancher-hardened-kubernetes-compat-fips",
			Input:  "FIPS-compliant Kubernetes components compatible with Rancher's RKE2",
			Expect: true,
		},
		{
			Name:   "opentelemetry-nodejs-instrumentation-fips",
			Input:  "Auto instrumention for any Node application to capture telemetry from a number of popular libraries and frameworks.",
			Expect: true,
		},
		{
			Name:   "nginx-iamguarded",
			Input:  "Minimal Wolfi-based nginx HTTP, reverse proxy, mail proxy, and a generic TCP/UDP proxy server",
			Expect: true,
		},
		{
			Name:   "crossplane-aws-fips",
			Input:  "FIPS-compliant Crossplane providers for managing Amazon Web Services (AWS) services on Kubernetes.",
			Expect: true,
		},
		{
			Name:   "heartbeat",
			Input:  "Heartbeat periodically checks the status of your services and determine whether they are available.",
			Expect: true,
		},
		{
			Name:   "nextcloud-server-fips",
			Input:  "Nextcloud server, a safe home for all your data",
			Expect: true,
		},
		{
			Name:   "request-6274",
			Input:  "Minimal container image with Go programming language using the go-slim package",
			Expect: true,
		},
		{
			Name:   "apache-apisix",
			Input:  "Apache APISIX is a dynamic, real-time, high-performance API Gateway.",
			Expect: true,
		},
		{
			Name:   "envoy-gateway",
			Input:  "Manages Envoy Proxy as a Standalone or Kubernetes-based Application Gateway.",
			Expect: true,
		},
		{
			Name:   "traefik",
			Input:  "[Traefik](https://github.com/traefik/traefik) is a cloud native application proxy.",
			Expect: true,
		},
		{
			Name:   "mongodb-iamguarded-fips",
			Input:  "[MongoDB](https://www.mongodb.com/) is a document-oriented database management system. MongoDB is a popular example of a NoSQL database, and stores data in JSON-like documents.",
			Expect: true,
		},
		{
			Name:   "prometheus-process-exporter",
			Input:  "process-exporter is an agent that collects process-specific metrics from a system and exposes them in a format that can be ingested by Prometheus",
			Expect: true,
		},
		{
			Name:   "prometheus-adapter",
			Input:  "[prometheus-adapter](https://github.com/kubernetes-sigs/prometheus-adapter) is a Prometheus project used to collect Prometheus metrics in Kubernetes.",
			Expect: true,
		},
		{
			Name:   "request-5291",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "rails-fips",
			Input:  "Ruby on Rails (often just called \"Rails\") is a web-application framework that includes everything needed to create database-backed web applications according to the Model-View-Controller (MVC) pattern.",
			Expect: true,
		},
		{
			Name:   "podinfo-fips",
			Input:  "Podinfo is a tiny web application that provides Go microsrvice template for kubernetes",
			Expect: true,
		},
		{
			Name:   "yunikorn-web-fips",
			Input:  "Apache YuniKorn Web UI",
			Expect: true,
		},
		{
			Name:   "kube-rbac-proxy-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "rabbitmq-default-user-credential-updater",
			Input:  "Minimal image with [default-user-credential-updater](https://github.com/rabbitmq/default-user-credential-updater)",
			Expect: true,
		},
		{
			Name:   "local-volume-provisioner",
			Input:  "Static provisioner of local volumes",
			Expect: true,
		},
		{
			Name:   "openfga",
			Input:  "A high performance and flexible authorization/permission engine built for developers and inspired by Google Zanzibar",
			Expect: true,
		},
		{
			Name:   "metrics-server-fips",
			Input:  "Metrics Server is a Kubernetes component that collects and provides resource usage metrics (CPU, memory) for nodes and pods.",
			Expect: true,
		},
		{
			Name:   "thanos-operator-fips",
			Input:  "Minimal FIPS image with the [thanos-operator](https://github.com/banzaicloud/thanos-operator) for managing Thanos components in Kubernetes.",
			Expect: true,
		},
		{
			Name:   "go-ipfs",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "vector-fips",
			Input:  "Minimal FIPS compliant image with [Vector](https://vector.dev/), an end-to-end data observability pipeline",
			Expect: true,
		},
		{
			Name:   "emissary",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "crossplane",
			Input:  "[Crossplane](https://www.crossplane.io/) lets you build a control plane with Kubernetes-style declarative and API-driven configuration and management for anything.",
			Expect: true,
		},
		{
			Name:   "victoria-metrics-fips",
			Input:  "VictoriaMetrics standalone image is a fast, cost-effective and scalable monitoring solution and time series database",
			Expect: true,
		},
		{
			Name:   "rancher-security-scan-fips",
			Input:  "Evaluates Kubernetes cluster security posture against established best practices using kube-bench framework.",
			Expect: true,
		},
		{
			Name:   "memcached",
			Input:  "[Memcached](https://memcached.org/) is an in-memory key-value store for small chunks of arbitrary data (strings, objects) from results of database calls, API calls, or page rendering.",
			Expect: true,
		},
		{
			Name:   "powershell",
			Input:  "Minimal Wolfi image with Powershell",
			Expect: true,
		},
		{
			Name:   "dbgate",
			Input:  "DbGate is a database administration tool for SQL Server, MySQL, PostgreSQL, MongoDB, Redis and SQLite. DbGate provides a modern web interface for database management, query execution, and data visualization.",
			Expect: true,
		},
		{
			Name:   "envoy-iamguarded",
			Input:  "[Envoy](https://www.envoyproxy.io/) Cloud-native high-performance edge/middle/service proxy",
			Expect: true,
		},
		{
			Name:   "cluster-api-provider-vsphere-fips",
			Input:  "Kubernetes-native declarative infrastructure for vSphere",
			Expect: true,
		},
		{
			Name:   "kube-bench",
			Input:  "Minimal image with [kube-bench](https://github.com/aquasecurity/kube-bench).",
			Expect: true,
		},
		{
			Name:   "chromium",
			Input:  "Minimal [Chromium](https://chromium.googlesource.com/chromium/src/) container image.",
			Expect: true,
		},
		{
			Name:   "plugin-barman-cloud",
			Input:  "CloudNativePG barman-cloud plugin for PostgreSQL backup and recovery to S3-compatible storage",
			Expect: true,
		},
		{
			Name:   "redis-operator-fips",
			Input:  "A FIPS-compliant redis-operator image that automates Redis cluster deployment, scaling, and management in Kubernetes environments.",
			Expect: true,
		},
		{
			Name:   "cloud-provider-azure-controller-manager",
			Input:  "Controller manager for Azure CLI",
			Expect: true,
		},
		{
			Name:   "argocd-iamguarded-fips",
			Input:  "argocd-iamguarded-fips is the FIPS-enabled IAMGuarded version of [Argo CD](https://argo-cd.readthedocs.io/en/stable/), a declarative continuous deployment tool for Kubernetes.",
			Expect: true,
		},
		{
			Name:   "step-ca",
			Input:  "Minimal image of [step-ca](https://smallstep.com/docs/step-ca), an online Certificate Authority (CA) for secure, automated X.509 and SSH certificate management",
			Expect: true,
		},
		{
			Name:   "rabbitmq-iamguarded",
			Input:  "[RabbitMQ](https://github.com/rabbitmq/rabbitmq-server) RabbitMQ is a message broker.",
			Expect: true,
		},
		{
			Name:   "contour-fips",
			Input:  "Contour is an ingress controller for Kubernetes that works by deploying the Envoy proxy as a reverse proxy and load balancer. Contour supports dynamic configuration updates out of the box while maintaining a lightweight profile.",
			Expect: true,
		},
		{
			Name:   "kubernetes-dns-node-cache-fips",
			Input:  "Minimal image that acts as a drop-in replacement for the [NodeLocal DNSCache](https://github.com/kubernetes/dns) image.",
			Expect: true,
		},
		{
			Name:   "psqlodbc-fips",
			Input:  "This image contains psqlodbc drivers for use with unixODBC.",
			Expect: true,
		},
		{
			Name:   "keycloak-config-cli",
			Input:  "Import YAML/JSON-formatted configuration files into Keycloak - Configuration as Code for Keycloak.",
			Expect: true,
		},
		{
			Name:   "argocd-iamguarded",
			Input:  "argocd-iamguarded is the IAMGuarded version of [Argo CD](https://argo-cd.readthedocs.io/en/stable/), a declarative continuous deployment tool for Kubernetes.",
			Expect: true,
		},
		{
			Name:   "mongodb-kubernetes-operator",
			Input:  "Chainguard's MongoDB Kubernetes Operator image enables you to deploy a MongoDB community instance to a Kubernetes cluster, as well as support replica sets, scaling the replicas up or down, version upgrades, custom roles, and TLS security.",
			Expect: true,
		},
		{
			Name:   "valkey-iamguarded-fips",
			Input:  "Valkey is an open source, in-memory data store used by millions of developers as a cache, vector database, document database, streaming engine, and message broker.",
			Expect: true,
		},
		{
			Name:   "spegel-fips",
			Input:  "Stateless cluster local OCI registry mirror.",
			Expect: true,
		},
		{
			Name:   "k8sgpt",
			Input:  "Minimal [k8sgpt](https://k8sgpt.ai/) container image.",
			Expect: true,
		},
		{
			Name:   "ratify-fips",
			Input:  "Artifact Ratification Framework (CNCF Sandbox)",
			Expect: true,
		},
		{
			Name:   "k8s-agents-operator",
			Input:  "k8s-agents-operator auto-instruments containerized workloads in Kubernetes with New Relic agents.",
			Expect: true,
		},
		{
			Name:   "pypy",
			Input:  "PyPy is a fast and compliant implementation of the Python language.",
			Expect: true,
		},
		{
			Name:   "pgpool2_exporter",
			Input:  "Prometheus exporter image for Pgpool-II metrics.",
			Expect: true,
		},
		{
			Name:   "request-5663",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "gogatekeeper",
			Input:  "Minimalist Wolfi-based image of Gatekeeper, an OpenID / Proxy service.",
			Expect: true,
		},
		{
			Name:   "min-toolkit-debug",
			Input:  "Wolfi container image with some debugging utilities included. Suitable for using as a debugging tool.",
			Expect: true,
		},
		{
			Name:   "kyverno-policy-reporter-plugin-kyverno",
			Input:  "This Plugin for Policy Reporter brings additional Kyverno specific information to the Policy Reporter UI",
			Expect: true,
		},
		{
			Name:   "kafbat-ui",
			Input:  "[Kafbat-UI](https://kafbat.io/) is an open-source Web UI for managing Apache Kafka clusters",
			Expect: true,
		},
		{
			Name:   "altinity-clickhouse-server",
			Input:  "Altinity Stable Build for ClickHouse",
			Expect: true,
		},
		{
			Name:   "openscap",
			Input:  "NIST Certified SCAP 1.2 toolkit",
			Expect: true,
		},
		{
			Name:   "k8ssandra-client",
			Input:  "A kubectl plugin to simplify usage of k8ssandra",
			Expect: true,
		},
		{
			Name:   "prometheus-yet-another-cloudwatch-exporter",
			Input:  "Prometheus exporter for AWS CloudWatch.",
			Expect: true,
		},
		{
			Name:   "git",
			Input:  "A minimal Git image.",
			Expect: true,
		},
		{
			Name:   "nginx-s3-gateway",
			Input:  "NGINX S3 Gateway",
			Expect: true,
		},
		{
			Name:   "glibc-openssl",
			Input:  "The GNU C Library (glibc) is a C standard library implementation maintained by the GNU Project. This container image also contains OpenSSL, a software library for applications providing secure communications over a network.",
			Expect: true,
		},
		{
			Name:   "redpanda",
			Input:  "Redpanda is a Kafka-compatible streaming data platform with no JVM dependencies.",
			Expect: true,
		},
		{
			Name:   "gcc-glibc",
			Input:  "Minimal GCC (GNU Compiler Collection) image for building C applications compatible with glibc.",
			Expect: true,
		},
		{
			Name:   "lvm-driver",
			Input:  "Dynamically provision Stateful Persistent Node-Local Volumes & Filesystems for Kubernetes that is integrated with a backend LVM2 data storage stack.",
			Expect: true,
		},
		{
			Name:   "cert-manager-webhook-pdns",
			Input:  "A PowerDNS webhook for cert-manager",
			Expect: true,
		},
		{
			Name:   "curl",
			Input:  "Minimal [curl](https://curl.se/) image base containing curl and ca-certificates.",
			Expect: true,
		},
		{
			Name:   "airflow-core",
			Input:  "Apache Airflow offers a platform to author, schedule, and monitor workflows programmatically. This image is a minimal, slimmed-down version of the official Apache Airflow with only core components.",
			Expect: true,
		},
		{
			Name:   "karma",
			Input:  "A dashboard for managing alerts from Alertmanager",
			Expect: true,
		},
		{
			Name:   "opentelemetry-operator",
			Input:  "Kubernetes Operator for OpenTelemetry Collector",
			Expect: true,
		},
		{
			Name:   "tesseract",
			Input:  "Minimal image that contains tesseract",
			Expect: true,
		},
		{
			Name:   "ksops",
			Input:  "KSOPS, or kustomize-SOPS, is a kustomize KRM exec plugin for SOPS encrypted resources. KSOPS can be used to decrypt any Kubernetes resource, but is most commonly used to decrypt encryptedKubernetes Secrets and ConfigMaps. As a kustomize plugin, KSOPS allows you to manage, build, and apply encrypted manifests the same way you manage the rest of your Kubernetes manifests.",
			Expect: true,
		},
		{
			Name:   "postgres-repmgr-iamguarded-fips",
			Input:  "PostgreSQL HA, a cluster solution using the PostgreSQL replication manager.",
			Expect: true,
		},
		{
			Name:   "k8s_gateway-fips",
			Input:  "A CoreDNS plugin to resolve all types of external Kubernetes resources",
			Expect: true,
		},
		{
			Name:   "pvc-autoresizer",
			Input:  "pvc-autoresizer is a Kubernetes controller that monitors persistent volume claims (PVCs) and automatically resizes them based on usage metrics collected from Prometheus.",
			Expect: true,
		},
		{
			Name:   "ingress-nginx-controller-iamguarded-fips",
			Input:  " Ingress-NGINX Controller for Kubernetes",
			Expect: true,
		},
		{
			Name:   "temporal-admin-tools-fips",
			Input:  "Administrative command-line tools for Temporal workflow management",
			Expect: true,
		},
		{
			Name:   "crossplane-fips",
			Input:  "The Cloud Native Control Plane",
			Expect: true,
		},
		{
			Name:   "blob-csi-fips",
			Input:  "This driver allows Kubernetes to access Azure Storage via azure-storage-fuse & NFSv3.",
			Expect: true,
		},
		{
			Name:   "busybox",
			Input:  "Container image with only busybox and libc (available in both musl and glibc variants). Suitable for running any binaries that only have a dependency on glibc/musl.",
			Expect: true,
		},
		{
			Name:   "zot",
			Input:  "Minimal image with [zot](https://github.com/project-zot/zot) binary.",
			Expect: true,
		},
		{
			Name:   "apache-nifi",
			Input:  "Apache NiFi was made for dataflow. It supports highly configurable directed graphs of data routing, transformation, and system mediation logic.",
			Expect: true,
		},
		{
			Name:   "consul-k8s-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "statsd",
			Input:  "Daemon for easy but powerful stats aggregation",
			Expect: true,
		},
		{
			Name:   "vela-core",
			Input:  "KubeVela is a modern software delivery platform that makes deploying and operating applications across today's hybrid, multi-cloud environments easier, faster and more reliable.",
			Expect: true,
		},
		{
			Name:   "rabbitmq-cluster-operator",
			Input:  "RabbitMQ Cluster Kubernetes Operator",
			Expect: true,
		},
		{
			Name:   "wait-for-it",
			Input:  "Container image for testing whether a service is listening on an address/port combination.",
			Expect: true,
		},
		{
			Name:   "nextcloud-server",
			Input:  "Nextcloud server, a safe home for all your data",
			Expect: true,
		},
		{
			Name:   "prometheus-pgbouncer-exporter-fips",
			Input:  "A FIPS-compliant image for PgBouncer Exporter. A Prometheus exporter that collects and exposes metrics from PgBouncer, a lightweight connection pooler for PostgreSQL",
			Expect: true,
		},
		{
			Name:   "atlantis",
			Input:  "Terraform Pull Request Automation",
			Expect: true,
		},
		{
			Name:   "bank-vaults",
			Input:  "Minimal Image for [Bank Vaults](https://bank-vaults.dev/), a CLI tool to init, unseal and configure Vault ",
			Expect: true,
		},
		{
			Name:   "gitea",
			Input:  "A painless self-hosted all-in-one software development service, including Git hosting, code review, team collaboration, package registry and CI/CD",
			Expect: true,
		},
		{
			Name:   "nats",
			Input:  "NATS is a flexible messaging system providing pub/sub, streaming, storage etc.",
			Expect: true,
		},
		{
			Name:   "vllm-openai-fips",
			Input:  "vLLM is a high-throughput and memory-efficient inference engine for Large Language Models (LLMs). This FIPS-validated variant provides OpenSSL FIPS 140-3 compliance for secure, production LLM deployments.",
			Expect: true,
		},
		{
			Name:   "gitness",
			Input:  "Minimal image with the `gitness` [server application](https://github.com/harness/gitness).",
			Expect: true,
		},
		{
			Name:   "aws-cli-iamguarded",
			Input:  "Minimal IAMGuarded image with the [AWS CLI](https://aws.amazon.com/cli/).",
			Expect: true,
		},
		{
			Name:   "authentik-fips",
			Input:  "[Authentik](https://goauthentik.io/) is an open-source Identity Provider that provides single sign-on with support for SAML, OAuth2/OIDC, LDAP, and RADIUS protocols.",
			Expect: true,
		},
		{
			Name:   "crossplane-keycloak",
			Input:  "A Crossplane provider for Keycloak.",
			Expect: true,
		},
		{
			Name:   "k8s-wait-for",
			Input:  "Container image for waiting for a k8s service, job or pods to enter a desired state.",
			Expect: true,
		},
		{
			Name:   "nova-fips",
			Input:  "FIPS-compliant version of Nova. Scans your Kubernetes cluster for outdated Helm charts, then suggests updates.",
			Expect: true,
		},
		{
			Name:   "request-1862",
			Input:  "The Triton Inference Server provides an optimized cloud and edge inferencing solution.",
			Expect: true,
		},
		{
			Name:   "ml-metadata-store-server",
			Input:  "[ML Metadata (MLMD)](https://www.tensorflow.org/tfx/guide/mlmd) remote gRPC server",
			Expect: true,
		},
		{
			Name:   "flannel-fips",
			Input:  "A FIPS-compliant Flannel image for Kubernetes. FIPS-compliant Flannel is a network fabric for Kubernetes, providing a way to manage network configurations across a cluster.",
			Expect: true,
		},
		{
			Name:   "sql_exporter",
			Input:  "Database-agnostic SQL Exporter for Prometheus",
			Expect: true,
		},
		{
			Name:   "kubernetes-csi-livenessprobe",
			Input:  " A sidecar container that can be included in a CSI plugin pod to enable integration with Kubernetes Liveness Probe.",
			Expect: true,
		},
		{
			Name:   "images/opencv",
			Input:  "OpenCV is a C++ computer vision and machine learning software library for image/video processing, object detection, face recognition, and more.",
			Expect: true,
		},
		{
			Name:   "node-problem-detector-fips",
			Input:  "[Node-problem-detector](https://github.com/kubernetes/node-problem-detector) aims to make various node problems visible to the upstream layers in the cluster management stack.",
			Expect: true,
		},
		{
			Name:   "rancher-fleet",
			Input:  "Deploy workloads from Git to large fleets of Kubernetes clusters",
			Expect: true,
		},
		{
			Name:   "splunk-otel-collector-fips",
			Input:  "Splunk OpenTelemetry Collector is a distribution of the OpenTelemetry Collector. It provides a unified way to receive, process, and export metric, trace, and log data for Splunk Observability Cloud",
			Expect: true,
		},
		{
			Name:   "minio-object-browser-iamguarded-fips",
			Input:  "MinIO Console is a library that provides a management and browser UI overlay for the MinIO Server",
			Expect: true,
		},
		{
			Name:   "postgres",
			Input:  "Minimal image for PostgreSQL, an advanced object-relational database management system.",
			Expect: true,
		},
		{
			Name:   "code-server",
			Input:  "VS Code in the browser",
			Expect: true,
		},
		{
			Name:   "livekit-egress",
			Input:  "livekit-egress is an open-source media egress service for real-time audio, video, and data, designed for low latency and scalability",
			Expect: true,
		},
		{
			Name:   "kubectl-iamguarded-fips",
			Input:  "Minimal image with kubectl binary.",
			Expect: true,
		},
		{
			Name:   "rabbitmq-cluster-operator-iamguarded",
			Input:  "RabbitMQ Cluster Kubernetes Operator",
			Expect: true,
		},
		{
			Name:   "ruby",
			Input:  "Minimal Ruby base image.",
			Expect: true,
		},
		{
			Name:   "nats-server-config-reloader-fips",
			Input:  "Monitors NATS configuration files and triggers reloads without restarting the server.",
			Expect: true,
		},
		{
			Name:   "paranoia",
			Input:  "Minimalist Wolfi-based paranoia image for inspecting certificate authorities in container images",
			Expect: true,
		},
		{
			Name:   "sigstore-policy-controller",
			Input:  "Policy Controller image that is part of the Sigstore stack",
			Expect: true,
		},
		{
			Name:   "text-generation-inference",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "knative-eventing-fips",
			Input:  "Event-driven application platform for Kubernetes",
			Expect: true,
		},
		{
			Name:   "postgis",
			Input:  "PostGIS extends the capabilities of the PostgreSQL relational database by adding support for storing, indexing, and querying geospatial data",
			Expect: true,
		},
		{
			Name:   "kubernetes-reflector",
			Input:  "Kubernetes controller for reflecting ConfigMaps, Secrets, and Certificates across namespaces",
			Expect: true,
		},
		{
			Name:   "generic-device-plugin",
			Input:  "A Kubernetes device plugin to schedule generic Linux devices",
			Expect: true,
		},
		{
			Name:   "wasmer",
			Input:  "This image contains the `wasmer` tool which can be used to compile or run wasm binaries.",
			Expect: true,
		},
		{
			Name:   "trufflehog",
			Input:  "TruffleHog is a tool that allows you to discover, classify, validate, and analyze leaked credentials.",
			Expect: true,
		},
		{
			Name:   "vitess-lite",
			Input:  "Vitess is a database clustering system for horizontal scaling of MySQL through generalized sharding.",
			Expect: true,
		},
		{
			Name:   "logstash-iamguarded",
			Input:  "Logstash dynamically ingests, transforms, and ships your data regardless of format or complexity.",
			Expect: true,
		},
		{
			Name:   "seaweedfs-fips",
			Input:  "SeaweedFS is a fast distributed storage system for blobs, objects, files, and data lake, providing S3-compatible API and filesystem interface.",
			Expect: true,
		},
		{
			Name:   "tekton-chains-fips",
			Input:  "Tekton Chains is a Kubernetes Custom Resource Definition (CRD) controller that allows you to manage your supply chain security in Tekton.",
			Expect: true,
		},
		{
			Name:   "tileserver-gl-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "ceph-csi-operator-fips",
			Input:  "Operator for Ceph CSI driver management in Kubernetes (FIPS)",
			Expect: true,
		},
		{
			Name:   "tritonserver-pytorch-backend",
			Input:  "The Triton backend for the PyTorch TorchScript models.",
			Expect: true,
		},
		{
			Name:   "grpc-health-probe",
			Input:  "A tool to perform health-checks for gRPC applications in Kubernetes and elsewhere",
			Expect: true,
		},
		{
			Name:   "pvc-autoresizer-fips",
			Input:  "pvc-autoresizer is a Kubernetes controller that monitors persistent volume claims (PVCs) and automatically resizes them based on usage metrics collected from Prometheus.",
			Expect: true,
		},
		{
			Name:   "celeborn",
			Input:  "Celeborn is dedicated to improving the efficiency and elasticity of different map-reduce engines and provides an elastic, high-efficient management service for intermediate data including shuffle data, spilled data, result data, etc. Currently, Celeborn is focusing on shuffle data.",
			Expect: true,
		},
		{
			Name:   "consul",
			Input:  "Minimal image with [Consul](https://www.consul.io/).",
			Expect: true,
		},
		{
			Name:   "filebeat-fips",
			Input:  "[filebeat](https://github.com/elastic/beats/tree/main/filebeat) Tails and ships log files",
			Expect: true,
		},
		{
			Name:   "redis-sentinel",
			Input:  "Minimal [redis-sentinel](https://redis.io/docs/latest/operate/oss_and_stack/management/sentinel/) Image which is compatible with [Bitnami's](https://github.com/bitnami/containers/tree/main/bitnami/redis-sentinel)",
			Expect: true,
		},
		{
			Name:   "grafana-rollout-operator",
			Input:  "Kubernetes Rollout Operator coordinates the rollout of pods between different StatefulSets within a specific namespace, and can be used to manage multi-AZ deployments",
			Expect: true,
		},
		{
			Name:   "mesosphere-vsphere-csi-syncer",
			Input:  "vSphere storage Container Storage Interface (CSI) plugin",
			Expect: true,
		},
		{
			Name:   "kubelet-csr-approver",
			Input:  "Kubernetes controller to enable automatic kubelet CSR validation after a series of (configurable) security checks",
			Expect: true,
		},
		{
			Name:   "regclient",
			Input:  "regclient is a client interface to OCI conformant registries and content shipped with the OCI Image Layout",
			Expect: true,
		},
		{
			Name:   "kubernetes-csi-external-attacher-fips",
			Input:  "Chainguard image for [kubernetes-csi/external-attacher](https://github.com/kubernetes-csi/external-attacher). Watches Kubernetes VolumeAttachment objects and triggers ControllerPublish/Unpublish against a CSI endpoint.",
			Expect: true,
		},
		{
			Name:   "zipkin",
			Input:  "Zipkin is a distributed tracing system.",
			Expect: true,
		},
		{
			Name:   "httpd",
			Input:  "[httpd](https://github.com/apache/httpd), a powerful and flexible HTTP/1.1 compliant web server.",
			Expect: true,
		},
		{
			Name:   "psqlodbc",
			Input:  "This image contains psqlodbc drivers for use with unixODBC.",
			Expect: true,
		},
		{
			Name:   "meilisearch",
			Input:  "Minimal meilisearch image.",
			Expect: true,
		},
		{
			Name:   "rancher-system-upgrade-controller",
			Input:  "A general-purpose, Kubernetes-native upgrade controller for nodes that provides automated upgrade capabilities for Kubernetes clusters.",
			Expect: true,
		},
		{
			Name:   "polaris",
			Input:  "Polaris is an open source policy engine for Kubernetes that validates and remediates resource configuration.",
			Expect: true,
		},
		{
			Name:   "nginx",
			Input:  "Minimal Wolfi-based nginx HTTP, reverse proxy, mail proxy, and a generic TCP/UDP proxy server",
			Expect: true,
		},
		{
			Name:   "druid",
			Input:  "Apache Druid is a high performance real-time analytics database.",
			Expect: true,
		},
		{
			Name:   "tekton-fips",
			Input:  "[Tekton](https://tekton.dev) provides a cloud-native Pipeline resource, mainly intended for CI/CD use cases.",
			Expect: true,
		},
		{
			Name:   "k8s-mig-manager",
			Input:  "MIG Partition Editor for NVIDIA GPUs",
			Expect: true,
		},
		{
			Name:   "cassandra",
			Input:  "[Cassandra](https://cassandra.apache.org) is a free and open-source, distributed, wide-column store, NoSQL database.",
			Expect: true,
		},
		{
			Name:   "crossplane-gcp-fips",
			Input:  "FIPS-compliant Crossplane GCP Providers deliver Kubernetes-native APIs for provisioning and managing Google Cloud resources through Crossplane.",
			Expect: true,
		},
		{
			Name:   "kubelet-csr-approver-fips",
			Input:  "Kubernetes controller to enable automatic kubelet CSR validation after a series of (configurable) security checks",
			Expect: true,
		},
		{
			Name:   "kube-logging-operator-node-exporter-fips",
			Input:  "Custom runner based Prometheus node exporter for kube-logging logging-operator",
			Expect: true,
		},
		{
			Name:   "aws-network-policy-agent-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "kratos",
			Input:  "A microservice-oriented governance framework written in Go",
			Expect: true,
		},
		{
			Name:   "emqx-exporter",
			Input:  "The emqx-exporter is designed to expose partial metrics that are not included in the EMQX Prometheus API.",
			Expect: true,
		},
		{
			Name:   "pgpool2-iamguarded",
			Input:  "open-source middleware that operates between PostgreSQL servers and clients, providing features such as connection pooling, load balancing, and replication to enhance database performance and availability",
			Expect: true,
		},
		{
			Name:   "opencost-fips",
			Input:  "OpenCost give teams visibility into current and historical Kubernetes and cloud spend and resource allocation.",
			Expect: true,
		},
		{
			Name:   "aws-efs-csi-driver",
			Input:  "Minimal images for [aws-efs-csi-driver](https://aws.amazon.com/efs/).",
			Expect: true,
		},
		{
			Name:   "memcached-fips",
			Input:  "[Memcached](https://memcached.org/) is an in-memory key-value store for small chunks of arbitrary data (strings, objects) from results of database calls, API calls, or page rendering.",
			Expect: true,
		},
		{
			Name:   "prometheus-cloudwatch-exporter",
			Input:  "Minimalist Wolfi-based Prometheus CloudWatch Exporter image for exporting metrics to Amazon AWS CloudWatch.",
			Expect: true,
		},
		{
			Name:   "tiktoken",
			Input:  "tiktoken is a fast BPE tokeniser for use with OpenAI's models",
			Expect: true,
		},
		{
			Name:   "timestamp-authority",
			Input:  "timestamp-authority is an RFC3161 Timestamp Authority, a core component of the sigstore stack",
			Expect: true,
		},
		{
			Name:   "apache-apisix",
			Input:  "Apache APISIX is a dynamic, real-time, high-performance API Gateway.",
			Expect: true,
		},
		{
			Name:   "clickhouse-iamguarded",
			Input:  "Minimal Wolfi-based ClickHouse analytics database image. [Clickhouse](https://clickhouse.com) is the fastest and most resource efficient open-source database for real-time apps and analytics.",
			Expect: true,
		},
		{
			Name:   "etcd-iamguarded-fips",
			Input:  "[etcd](https://etcd.io/) Distributed reliable key-value store for the most critical data of a distributed system",
			Expect: true,
		},
		{
			Name:   "prometheus-redis-exporter-fips",
			Input:  "Minimalist Wolfi-based Prometheus Redis Exporter image for exporting metrics to Redis.",
			Expect: true,
		},
		{
			Name:   "chisel",
			Input:  "A fast TCP/UDP tunnel over HTTP",
			Expect: true,
		},
		{
			Name:   "configmap-reload",
			Input:  "`configmap-reload` is a simple binary to trigger a reload when Kubernetes ConfigMaps or Secrets, mounted into pods, are updated.",
			Expect: true,
		},
		{
			Name:   "linkerd",
			Input:  "Ultralight, security-first service mesh for Kubernetes. Main repo for Linkerd 2.x.",
			Expect: true,
		},
		{
			Name:   "label-studio",
			Input:  "[Label Studio](https://labelstud.io/) is an open-source data labeling platform that supports annotation of audio, text, images, videos, and time series data.",
			Expect: true,
		},
		{
			Name:   "prometheus-process-exporter-fips",
			Input:  "process-exporter is an agent that collects process-specific metrics from a system and exposes them in a format that can be ingested by Prometheus",
			Expect: true,
		},
		{
			Name:   "crossplane-provider-gitlab-fips",
			Input:  "This image contains the Crossplane GitLab provider, which allows you to manage GitLab resources using Crossplane.",
			Expect: true,
		},
		{
			Name:   "unbound-mailcow",
			Input:  "Unbound is a validating, recursive, and caching DNS resolver.",
			Expect: true,
		},
		{
			Name:   "generic-device-plugin-fips",
			Input:  "A FIPS-compliant Kubernetes device plugin to schedule generic Linux devices",
			Expect: true,
		},
		{
			Name:   "request-7082",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "mysql",
			Input:  "MySQL is a widely used open-source relational database management system.",
			Expect: true,
		},
		{
			Name:   "jupyter-base-notebook",
			Input:  "Minimal Jupyter base notebook image using pip",
			Expect: true,
		},
		{
			Name:   "grafana-image-renderer",
			Input:  "A Grafana backend plugin that handles rendering of panels & dashboards to PNGs using headless browser (Chromium/Chrome)",
			Expect: true,
		},
		{
			Name:   "temporal-admin-tools",
			Input:  "Administrative command-line tools for Temporal workflow management",
			Expect: true,
		},
		{
			Name:   "static",
			Input:  "Base images with the minimum contents needed to run static binaries.",
			Expect: true,
		},
		{
			Name:   "pdns-recursor-fips",
			Input:  "PowerDNS Recursor is a non authoritative/recursing DNS server.",
			Expect: true,
		},
		{
			Name:   "tempo",
			Input:  "Grafana Tempo is a high volume, minimal dependency distributed tracing backend.",
			Expect: true,
		},
		{
			Name:   "apache-activemq-artemis",
			Input:  "ActiveMQ Artemis",
			Expect: true,
		},
		{
			Name:   "min-toolkit-debug-fips",
			Input:  "Wolfi container image with some debugging utilities included. Suitable for using as a debugging tool.",
			Expect: true,
		},
		{
			Name:   "prometheus-postgres-exporter",
			Input:  "A PostgreSQL metric exporter for Prometheus",
			Expect: true,
		},
		{
			Name:   "kubernetes-event-exporter",
			Input:  "Minimalist [wolfi](https://github.com/wolfi-dev)-based image of [Kubernetes Event Exporter](https://github.com/resmoio/kubernetes-event-exporter) for exporting Kubernetes events to various outputs to be used for observability or alerting purposes.",
			Expect: true,
		},
		{
			Name:   "tensorflow",
			Input:  "An Open Source Machine Learning Framework for Everyone",
			Expect: true,
		},
		{
			Name:   "dapr-fips",
			Input:  "Dapr is a portable, event-driven, runtime for building distributed applications across cloud and edge.",
			Expect: true,
		},
		{
			Name:   "kapp-controller",
			Input:  "Continuous delivery and package management for Kubernetes.",
			Expect: true,
		},
		{
			Name:   "request-6338",
			Input:  "Chainguard's `static` image, without `glibc-locale-posix`",
			Expect: true,
		},
		{
			Name:   "crossplane-provider-kubernetes",
			Input:  "This image contains the Crossplane Kubernetes provider, which allows you to manage Kubernetes resources using Crossplane.",
			Expect: true,
		},
		{
			Name:   "sqlite3",
			Input:  "SQLite is a C-language library that implements a small, fast, self-contained, high-reliability, full-featured, SQL database engine.",
			Expect: true,
		},
		{
			Name:   "clamav-fips",
			Input:  "ClamAV® is an open source antivirus engine for detecting trojans, viruses, malware & other malicious threats.",
			Expect: true,
		},
		{
			Name:   "spicedb-operator-fips",
			Input:  "This is the FIPS-compliant variant of the SpiceDB Operator, a Kubernetes operator for managing SpiceDB clusters, providing automated deployment, scaling, and management of SpiceDB instances.",
			Expect: true,
		},
		{
			Name:   "mlflow-iamguarded",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "apache-nifi-registry",
			Input:  "Registry for storing and managing shared resources such as versioned flows across one or more instances of NiFi.",
			Expect: true,
		},
		{
			Name:   "kube-metrics-adapter-fips",
			Input:  "Minimal Adapter to expose custom metrics to Kubernetes HPA via Prometheus",
			Expect: true,
		},
		{
			Name:   "rabbitmq-messaging-topology-operator",
			Input:  "RabbitMQ messaging topology operator",
			Expect: true,
		},
		{
			Name:   "dfc",
			Input:  "CLI to convert Dockerfiles to use Chainguard Images and APKs in FROM and RUN lines",
			Expect: true,
		},
		{
			Name:   "nvidia-container-toolkit",
			Input:  "The NVIDIA Container Toolkit allows users to build and run GPU accelerated containers.",
			Expect: true,
		},
		{
			Name:   "prometheus-alertmanager-iamguarded",
			Input:  "The Alertmanager handles alerts sent by client applications such as the Prometheus server",
			Expect: true,
		},
		{
			Name:   "images/hailo-ai-onnxruntime",
			Input:  "Container image with ONNX Runtime, HailoRT, and the Hailo Execution Provider for hardware-accelerated ML inference",
			Expect: true,
		},
		{
			Name:   "envoy-iamguarded-fips",
			Input:  "[Envoy](https://www.envoyproxy.io/) Cloud-native high-performance edge/middle/service proxy",
			Expect: true,
		},
		{
			Name:   "newrelic-fluent-bit-output",
			Input:  "Minimal [newrelic-fluent-bit-output](https://github.com/newrelic/newrelic-fluent-bit-output) container image.",
			Expect: true,
		},
		{
			Name:   "perl-fips",
			Input:  "Container image for building Perl applications with FIPS.",
			Expect: true,
		},
		{
			Name:   "configurable-http-proxy",
			Input:  "configurable-http-proxy provides you with a way to update and manage a proxy table using a command line interface or REST API",
			Expect: true,
		},
		{
			Name:   "flux-operator",
			Input:  "Flux Operator is a Kubernetes controller for managing the lifecycle of Flux CD",
			Expect: true,
		},
		{
			Name:   "kube-state-metrics-iamguarded-fips",
			Input:  "Kube-state-metrics generates Prometheus metrics about Kubernetes objects",
			Expect: true,
		},
		{
			Name:   "cert-manager-iamguarded",
			Input:  "[cert-manager](https://cert-manager.io) is a tool for provisioning and managing TLS certificates in Kubernetes.",
			Expect: true,
		},
		{
			Name:   "cilium-fips",
			Input:  "[Cilium](https://cilium.io/) is an open source, cloud native solution for providing, securing, and observing network connectivity between workloads using eBPF",
			Expect: true,
		},
		{
			Name:   "gpu-operator-validator",
			Input:  "Minimal [gpu-operator-validator](https://github.com/NVIDIA/gpu-operator) container image.",
			Expect: true,
		},
		{
			Name:   "milvus",
			Input:  "High-performance, cloud-native vector database built for scalable vector ANN search",
			Expect: true,
		},
		{
			Name:   "node-problem-detector",
			Input:  "[Node-problem-detector](https://github.com/kubernetes/node-problem-detector) aims to make various node problems visible to the upstream layers in the cluster management stack.",
			Expect: true,
		},
		{
			Name:   "zabbix-agent2",
			Input:  "Minimalist Wolfi-based Zabbix Agent 2 for monitoring hosts and sending metrics to Zabbix Server.",
			Expect: true,
		},
		{
			Name:   "gpu-feature-discovery",
			Input:  "Minimal [gpu-feature-discovery](https://github.com/NVIDIA/gpu-feature-discovery) container image.",
			Expect: true,
		},
		{
			Name:   "portieris",
			Input:  "A Kubernetes Admission Controller for verifying image trust.",
			Expect: true,
		},
		{
			Name:   "opentelemetry-collector-contrib",
			Input:  "Minimal image with [opentelemetry-collector-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib).",
			Expect: true,
		},
		{
			Name:   "xeol",
			Input:  "A scanner for end-of-life (EOL) software and dependencies in container images, filesystems, and SBOMs",
			Expect: true,
		},
		{
			Name:   "task",
			Input:  "Task is a task runner and build tool that aims to be simpler and easier to use than GNU Make",
			Expect: true,
		},
		{
			Name:   "pdns-auth",
			Input:  "PowerDNS Authoritative Server - High-performance DNS server with flexible backend support.",
			Expect: true,
		},
		{
			Name:   "filebrowser",
			Input:  "Filebrowser provides a file managing interface within a specified directory and it can be used to upload, delete, preview, rename and edit your files",
			Expect: true,
		},
		{
			Name:   "swift",
			Input:  "Container image for building Swift applications.",
			Expect: true,
		},
		{
			Name:   "fluent-bit-watcher",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "guacamole-server",
			Input:  "Minimal [Guacamole server](https://guacamole.apache.org/) remote desktop gateway image.",
			Expect: true,
		},
		{
			Name:   "images/asciinema",
			Input:  "asciinema is a CLI tool for recording and sharing terminal sessions. It records terminal sessions as lightweight text-based recordings (casts) that can be played back, converted between formats, and shared.",
			Expect: true,
		},
		{
			Name:   "kubescape-operator",
			Input:  "Kubescape-Operator is an in-cluster component of the Kubescape security platform that orchestrates security scanning and policy enforcement.",
			Expect: true,
		},
		{
			Name:   "step-issuer",
			Input:  "Minimal container image of [step-issuer](https://smallstep.com/docs/platform/), a certificate issuer for cert-manager using step certificates CA",
			Expect: true,
		},
		{
			Name:   "bento",
			Input:  "Bento is a high performance and resilient stream processor, able to connect various sources and sinks in a range of brokering patterns and perform hydration, enrichments, transformations and filters on payloads.",
			Expect: true,
		},
		{
			Name:   "spegel",
			Input:  "Stateless cluster local OCI registry mirror.",
			Expect: true,
		},
		{
			Name:   "metacontroller",
			Input:  "Minimal Metacontroller Image",
			Expect: true,
		},
		{
			Name:   "linkerd-cni-plugin",
			Input:  "Init container that sets up the iptables rules to forward traffic into the Linkerd2 sidecar proxy",
			Expect: true,
		},
		{
			Name:   "opa",
			Input:  "Open Policy Agent (OPA) is an open source, general-purpose policy engine..",
			Expect: true,
		},
		{
			Name:   "stunnel",
			Input:  "This image contains the CLI for the [stunnel](https://www.stunnel.org/) networking tool",
			Expect: true,
		},
		{
			Name:   "amazon-corretto-jdk",
			Input:  "Amazon Corretto is a no-cost, multi-platform, production-ready distribution of corresponding version of OpenJDK",
			Expect: true,
		},
		{
			Name:   "metallb",
			Input:  "[MetalLB](https://metallb.org) provides network load balancers for bare-metal Kubernetes clusters",
			Expect: true,
		},
		{
			Name:   "google-cloud-sdk-iamguarded",
			Input:  "Minimal IAMGuarded image with the [Google Cloud SDK](https://cloud.google.com/sdk/).",
			Expect: true,
		},
		{
			Name:   "trillian",
			Input:  "[Trillian](https://github.com/google/trillian) is a Merkle tree implementation that is used as the backing for various functionalities including Certificate Transparency and the Sigstore Rekor transparency log.",
			Expect: true,
		},
		{
			Name:   "haproxy-iamguarded",
			Input:  "A minimal [haproxy](https://www.haproxy.org/) base image rebuilt every night from source.",
			Expect: true,
		},
		{
			Name:   "buck2",
			Input:  "Minimal image with [buck2](https://buck2.build) build system binaries and toolchain.",
			Expect: true,
		},
		{
			Name:   "kube-vip-cloud-provider",
			Input:  "A general purpose cloud provider for kube-vip",
			Expect: true,
		},
		{
			Name:   "k8sgpt-operator",
			Input:  "Minimal k8sgpt-operator container image.",
			Expect: true,
		},
		{
			Name:   "qdrant",
			Input:  "[qdrant](https://github.com/qdrant/qdrant) Qdrant is a high-performance, massive-scale Vector Database for the next generation of AI.",
			Expect: true,
		},
		{
			Name:   "cosign",
			Input:  "Minimalist Wolfi-based Cosign images for signing and verifying images using Sigstore.",
			Expect: true,
		},
		{
			Name:   "bazel",
			Input:  "[Bazel](https://bazel.build) - A fast, scalable, multi-language and extensible build system.",
			Expect: true,
		},
		{
			Name:   "syft",
			Input:  "A tool for generating a Software Bill of Materials (SBOM) from container images and filesystems.",
			Expect: true,
		},
		{
			Name:   "gpu-operator-fips",
			Input:  "The NVIDIA GPU Operator bootstraps, configures, and manages GPUs in Kubernetes.",
			Expect: true,
		},
		{
			Name:   "shellcheck",
			Input:  "ShellCheck -- Shell script analysis tool",
			Expect: true,
		},
		{
			Name:   "request-6336",
			Input:  "Custom Base image on top of Chainguard glibc-openssl image (glibc base image bundled with OpenSSL)",
			Expect: true,
		},
		{
			Name:   "thanos-operator",
			Input:  "Minimal image with the [thanos-operator](https://github.com/banzaicloud/thanos-operator) for managing Thanos components in Kubernetes.",
			Expect: true,
		},
		{
			Name:   "kafka-iamguarded",
			Input:  "Apache Kafka distributed event store and stream-processing platform, compatible with iamguarded charts",
			Expect: true,
		},
		{
			Name:   "aws-cli",
			Input:  "Minimal [aws-cli](https://github.com/aws/aws-cli) container image.",
			Expect: true,
		},
		{
			Name:   "cluster-api-gcp-controller-fips",
			Input:  "Kubernetes Cluster API provider for Google Cloud Platform infrastructure management.",
			Expect: true,
		},
		{
			Name:   "povray",
			Input:  "POV-Ray is a ray-tracing program that generates images from text-based scene descriptions.",
			Expect: true,
		},
		{
			Name:   "tailscale",
			Input:  "[Tailscale](https://github.com/tailscale/tailscale) is a WireGuard-based mesh VPN",
			Expect: true,
		},
		{
			Name:   "dapr",
			Input:  "Dapr is a portable, event-driven, runtime for building distributed applications across cloud and edge.",
			Expect: true,
		},
		{
			Name:   "tigera-operator",
			Input:  "Minimal Project Calico Tigera Operator Image",
			Expect: true,
		},
		{
			Name:   "request-2354",
			Input:  "[Cilium](https://cilium.io/) is an open source, cloud native solution for providing, securing, and observing network connectivity between workloads using eBPF.",
			Expect: true,
		},
		{
			Name:   "prometheus-mongodb-exporter",
			Input:  "Minimalist Wolfi-based Prometheus MongoDB Exporter image for exporting various metrics about MongoDB.",
			Expect: true,
		},
		{
			Name:   "sbt",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "k3s",
			Input:  "Minimal image of [K3s](https://k3s.io/), a lightweight Kubernetes distribution",
			Expect: true,
		},
		{
			Name:   "falco",
			Input:  "A minimal, [wolfi](https://github.com/wolfi-dev)-based image for falco. This streamlined variant of [Falco](https://github.com/falcosecurity/falco/tree/master) designed for real-time security monitoring on Linux, replaces the traditional kernel module with eBPF technology, thus enhancing portability in containerized environments.",
			Expect: true,
		},
		{
			Name:   "nvidia-device-plugin",
			Input:  "Minimal [nvidia-device-plugin](https://github.com/NVIDIA/k8s-device-plugin) container image.",
			Expect: true,
		},
		{
			Name:   "ctlog",
			Input:  "ctlog is deployed as part of the sigstore stack",
			Expect: true,
		},
		{
			Name:   "openbao-fips",
			Input:  "Minimal image with OpenBao, FIPS compliant.",
			Expect: true,
		},
		{
			Name:   "k6-operator",
			Input:  "Kubernetes operator for running distributed k6 performance tests",
			Expect: true,
		},
		{
			Name:   "cadvisor",
			Input:  "[cAdvisor (Container Advisor)](https://github.com/google/cadvisor) provides container users an understanding of the resource usage and performance characteristics of their running containers.",
			Expect: true,
		},
		{
			Name:   "etcd",
			Input:  "[etcd](https://etcd.io/) Distributed reliable key-value store for the most critical data of a distributed system",
			Expect: true,
		},
		{
			Name:   "kubeflow-fips",
			Input:  "Kubeflow is a Machine Learning Toolkit for Kubernetes with FIPS",
			Expect: true,
		},
		{
			Name:   "argo-rollouts",
			Input:  "Argo Rollouts is a Kubernetes controller and set of CRDs which provide advanced deployment capabilities such as blue-green, canary, canary analysis, experimentation, and progressive delivery features to Kubernetes.",
			Expect: true,
		},
		{
			Name:   "thanos",
			Input:  "Minimal Thanos Image, a highly available Prometheus setup with long term storage",
			Expect: true,
		},
		{
			Name:   "rancher-webhook",
			Input:  "Rancher Webhook",
			Expect: true,
		},
		{
			Name:   "cloudflared",
			Input:  "Cloudflare Tunnel client (formerly Argo Tunnel)",
			Expect: true,
		},
		{
			Name:   "velero-plugin-for-microsoft-azure-fips",
			Input:  "Velero plugin for Microsoft Azure that provides backup and restore functionality for Azure Blob Storage and Azure Disk snapshots",
			Expect: true,
		},
		{
			Name:   "k6",
			Input:  "Load testing tool for testing APIs, microservices, and websites.",
			Expect: true,
		},
		{
			Name:   "litestream",
			Input:  "Container image for [litestream](https://litestream.io), to replicate SQLite databases.",
			Expect: true,
		},
		{
			Name:   "crossplane-provider-Kubernetes-fips",
			Input:  "This image contains the Crossplane Kubernetes provider, which allows you to manage Kubernetes resources using Crossplane.",
			Expect: true,
		},
		{
			Name:   "descheduler-fips",
			Input:  "Kubernetes Descheduler is a tool that evicts pods from nodes based on configurable policies to improve cluster balance, resource utilization, and scheduling efficiency.",
			Expect: true,
		},
		{
			Name:   "tekton",
			Input:  "[Tekton](https://tekton.dev) provides a cloud-native Pipeline resource, mainly intended for CI/CD use cases.",
			Expect: true,
		},
		{
			Name:   "crossplane-sql-fips",
			Input:  "A Crossplane provider for SQL FIPS version",
			Expect: true,
		},
		{
			Name:   "request-5774",
			Input:  "Minimal image with Apache Zookeeper 3.9.3 (Latest official release). ZooKeeper is a centralized service for maintaining configuration information, naming, providing distributed synchronization, and providing group services.",
			Expect: true,
		},
		{
			Name:   "az",
			Input:  "Azure CLI",
			Expect: true,
		},
		{
			Name:   "git-sync",
			Input:  "A sidecar app which clones a git repo and keeps it in sync with the upstream.",
			Expect: true,
		},
		{
			Name:   "rabbitmq-iamguarded-fips",
			Input:  "[RabbitMQ](https://github.com/rabbitmq/rabbitmq-server) RabbitMQ is a message broker.",
			Expect: true,
		},
		{
			Name:   "wavefront-proxy",
			Input:  "Minimal wavefront-proxy image",
			Expect: true,
		},
		{
			Name:   "opentofu",
			Input:  "[OpenTofu](https://opentofu.org/) is an open-source infrastructure as code tool that allows you to declaratively manage your cloud infrastructure. OpenTofu is a fork of Terraform managed by the Linux Foundation.",
			Expect: true,
		},
		{
			Name:   "verticadb-operator",
			Input:  "The VerticaDB operator automates tasks and monitors the state of your Vertica on Kubernetes deployments.",
			Expect: true,
		},
		{
			Name:   "zabbix-agent2-fips",
			Input:  "FIPS-validated Wolfi-based Zabbix Agent 2 for monitoring hosts and sending metrics to Zabbix Server.",
			Expect: true,
		},
		{
			Name:   "azurefile-csi",
			Input:  "This driver allows Kubernetes to access Azure File volume using smb and nfs protocols, csi plugin name: file.csi.azure.com.",
			Expect: true,
		},
		{
			Name:   "helm-operator",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "request-5548",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "cfssl-self-sign",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "sonobuoy",
			Input:  "A diagnostic tool for Kubernetes clusters to run conformance and plugin-based tests",
			Expect: true,
		},
		{
			Name:   "victoriametrics-operator",
			Input:  "Kubernetes operator for Victoria Metrics",
			Expect: true,
		},
		{
			Name:   "net-kourier",
			Input:  "Knative Ingress implementation using Envoy",
			Expect: true,
		},
		{
			Name:   "kayenta",
			Input:  "Automated Canary Service",
			Expect: true,
		},
		{
			Name:   "docker-dind-fips",
			Input:  "Chainguard image for Docker in Docker (DinD), allowing you to run Docker within a container.",
			Expect: true,
		},
		{
			Name:   "glibc",
			Input:  "The GNU C Library (glibc) is a C standard library implementation maintained by the GNU Project. glibc aims to provide a consistent interface to help developers write software that will work across multiple platforms.",
			Expect: true,
		},
		{
			Name:   "clickhouse-operator",
			Input:  "Kubernetes Operator for ClickHouse. Creates, configures and manages ClickHouse clusters running on Kubernetes.",
			Expect: true,
		},
		{
			Name:   "rancher-shell-fips",
			Input:  "Minimal FIPS compliant kubectl and helm installer image for Rancher",
			Expect: true,
		},
		{
			Name:   "kustomize-mutating-webhook-fips",
			Input:  "A dynamic solution to patch FluxCD Kustomization resources, seamlessly integrating and federating substitution variables across multiple namespaces.",
			Expect: true,
		},
		{
			Name:   "cluster-api-capd-manager",
			Input:  "Cluster API Provider Docker (CAPD) is a reference implementation of an infrastructure provider for the Cluster API project using Docker, not designed for production use and intended for development environments only.",
			Expect: true,
		},
		{
			Name:   "rancher",
			Input:  "Complete container management platform",
			Expect: true,
		},
		{
			Name:   "cass-operator",
			Input:  "[cass-operator](https://github.com/k8ssandra/cass-operator), is a Kubernetes operator for managing Apache Cassandra. It automates tasks like deployment, scaling, and configuration management, facilitating the integration of Cassandra clusters with Kubernetes environments.",
			Expect: true,
		},
		{
			Name:   "kubectl-iamguarded",
			Input:  "Minimal image with kubectl binary.",
			Expect: true,
		},
		{
			Name:   "db-operator-fips",
			Input:  "The DB Operator creates databases and make them available in the cluster via Custom Resource.",
			Expect: true,
		},
		{
			Name:   "crossplane-sql",
			Input:  "A Crossplane provider for SQL.",
			Expect: true,
		},
		{
			Name:   "terragrunt-fips",
			Input:  "A cloud infrastructure orchestration tool that supports OpenTofu/Terraform.",
			Expect: true,
		},
		{
			Name:   "duckdb",
			Input:  "DuckDB is an analytical in-process SQL database management system.",
			Expect: true,
		},
		{
			Name:   "cluster-autoscaler",
			Input:  "Minimal Kubernetes Cluster Autoscaler Image",
			Expect: true,
		},
		{
			Name:   "kafka-bridge-fips",
			Input:  "FIPS-compliant HTTP bridge for Apache Kafka using Vert.x framework",
			Expect: true,
		},
		{
			Name:   "prometheus-redis-exporter-iamguarded",
			Input:  "Prometheus Redis Exporter image for exporting metrics to Redis.",
			Expect: true,
		},
		{
			Name:   "kubernetes-replicator",
			Input:  "[kubernetes-replicator](https://github.com/mittwald/kubernetes-replicator) is a custom Kubernetes controller that can be used to make secrets and config maps available in multiple namespaces.",
			Expect: true,
		},
		{
			Name:   "prometheus-snmp-exporter",
			Input:  "SNMP Exporter for Prometheus",
			Expect: true,
		},
		{
			Name:   "plugin-barman-cloud-fips",
			Input:  "CloudNativePG barman-cloud plugin for PostgreSQL backup and recovery to S3-compatible storage",
			Expect: true,
		},
		{
			Name:   "tekton-cli",
			Input:  "The Tekton Pipelines CLI project provides a command-line interface (CLI) for interacting with Tekton, an open-source framework for Continuous Integration and Delivery (CI/CD) systems.",
			Expect: true,
		},
		{
			Name:   "newrelic-kube-events",
			Input:  "Minimal [newrelic-kube-events](https://github.com/newrelic/nri-kube-events) container image.",
			Expect: true,
		},
		{
			Name:   "chart-testing-fips",
			Input:  "Tool for testing Helm charts, used for linting and testing pull requests.",
			Expect: true,
		},
		{
			Name:   "calico",
			Input:  "[Calico](https://projectcalico.docs.tigera.io/) is a networking and security solution that enables Kubernetes workloads and non-Kubernetes/legacy workloads to communicate seamlessly and securely.",
			Expect: true,
		},
		{
			Name:   "git-iamguarded",
			Input:  "A minimal Git image for use with Iamguarded charts.",
			Expect: true,
		},
		{
			Name:   "gitlab-operator-fips",
			Input:  "Kubernetes Operator for GitLab Server",
			Expect: true,
		},
		{
			Name:   "custom-pod-autoscaler-operator-fips",
			Input:  "Operator for managing Kubernetes Custom Pod Autoscalers (CPA)",
			Expect: true,
		},
		{
			Name:   "k8s-secret-sync",
			Input:  "k8s-secret-sync provides two-way sync of JSON files to Kubernetes secret objects",
			Expect: true,
		},
		{
			Name:   "ko",
			Input:  "Minimal image to build and deploy Go applications using [ko](https://ko.build/)",
			Expect: true,
		},
		{
			Name:   "ant-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "step-cli",
			Input:  "[step-cli](https://smallstep.com/docs/step-cli) is an easy-to-use CLI tool for building, operating, and automating Public Key Infrastructure (PKI) systems and workflows",
			Expect: true,
		},
		{
			Name:   "camunda-keycloak",
			Input:  "Minimalist Wolfi-based [Camunda Keycloak](https://github.com/camunda/keycloak/) image for identity and access management.",
			Expect: true,
		},
		{
			Name:   "rabbitmq-default-user-credential-updater-iamguarded",
			Input:  "Image with [default-user-credential-updater](https://github.com/rabbitmq/default-user-credential-updater)",
			Expect: true,
		},
		{
			Name:   "postgres-iamguarded-fips",
			Input:  "PostgreSQL is a powerful, open source object-relational database system.",
			Expect: true,
		},
		{
			Name:   "request-5756",
			Input:  "GraalVM is an advanced JDK with ahead-of-time Native Image compilation",
			Expect: true,
		},
		{
			Name:   "k8s-secret-sync-fips",
			Input:  "k8s-secret-sync provides two-way sync of JSON files to Kubernetes secret objects",
			Expect: true,
		},
		{
			Name:   "ghidra",
			Input:  "Ghidra application image offering both GUI and headless modes of operation",
			Expect: true,
		},
		{
			Name:   "minio-iamguarded",
			Input:  "MinIO is a high-performance, S3 compatible object store. This iamguarded variant is specifically designed to work with the iamguarded Helm chart.",
			Expect: true,
		},
		{
			Name:   "cephcsi",
			Input:  "CephCSI is the Container Storage Interface (CSI) driver for Ceph, providing support for RBD and CephFS.",
			Expect: true,
		},
		{
			Name:   "linkerd-extension-init",
			Input:  "A utility for initializing Linkerd extension namespaces after installation",
			Expect: true,
		},
		{
			Name:   "flux-source-watcher",
			Input:  "Flux Source Watcher extends Flux CD with ArtifactGenerator CRD for source composition",
			Expect: true,
		},
		{
			Name:   "prometheus-blackbox-exporter",
			Input:  "Prometheus blackbox exporter allows blackbox probing of endpoints over HTTP, HTTPS, DNS, TCP, ICMP and gRPC.",
			Expect: true,
		},
		{
			Name:   "nvidia-gpu-driver",
			Input:  "Tools necessary for GPU and feature discovery for NVIDIA GPU driver container that allows the provisioning of the NVIDIA driver through the use of containers.",
			Expect: true,
		},
		{
			Name:   "images/secretgen-controller-fips",
			Input:  "secretgen-controller provides CRDs to specify what secrets need to be on Kubernetes cluster (to be generated or not)",
			Expect: true,
		},
		{
			Name:   "mattermost",
			Input:  "Community edition of Mattermost, a self-hostable chat service with file sharing, search, and integrations. It is designed as an internal chat for organisations and companies, and mostly markets itself as an open-source alternative to Slack and Microsoft Teams.",
			Expect: true,
		},
		{
			Name:   "promxy-fips",
			Input:  "Minimal image with Promxy, FIPS compliant.",
			Expect: true,
		},
		{
			Name:   "superset",
			Input:  "A minimal, [wolfi](https://github.com/wolfi-dev)-based image for Apache Superset. [Apache Superset](https://github.com/apache/superset/tree/master)is a Data Visualization and Data Exploration Platform",
			Expect: true,
		},
		{
			Name:   "opa-fips",
			Input:  "Open Policy Agent (OPA) is an open source, general-purpose policy engine..",
			Expect: true,
		},
		{
			Name:   "tetragon",
			Input:  "eBPF-based Security Observability and Runtime Enforcement",
			Expect: true,
		},
		{
			Name:   "memcached-exporter-iamguarded",
			Input:  "A memcached exporter for Prometheus.",
			Expect: true,
		},
		{
			Name:   "cockroach",
			Input:  "CockroachDB is a cloud-native distributed SQL database designed to build, scale, and manage modern, data-intensive applications.",
			Expect: true,
		},
		{
			Name:   "cluster-api-fips",
			Input:  "Home for Cluster API, a subproject of sig-cluster-lifecycle",
			Expect: true,
		},
		{
			Name:   "apm-server-fips",
			Input:  "Elastic APM is an application performance monitoring system built on the Elastic Stack.",
			Expect: true,
		},
		{
			Name:   "kuma",
			Input:  "The universal Envoy service mesh for distributed service connectivity",
			Expect: true,
		},
		{
			Name:   "livekit-server",
			Input:  "livekit-server is an open-source media server for real-time audio, video, and data, designed for low latency and scalability",
			Expect: true,
		},
		{
			Name:   "sftp",
			Input:  "Various scripts from https://github.com/atmoz/sftp to help run SFTP in a container",
			Expect: true,
		},
		{
			Name:   "rclone-fips",
			Input:  "Rclone syncs files and directories to and from different cloud storage providers.",
			Expect: true,
		},
		{
			Name:   "helm-operator-fips",
			Input:  "open source toolkit to manage Kubernetes native applications.",
			Expect: true,
		},
		{
			Name:   "cert-manager-istio-csr",
			Input:  "istio-csr is an agent that allows for Istio workload and control plane components to be secured using cert-manager.",
			Expect: true,
		},
		{
			Name:   "jaeger-iamguarded",
			Input:  "CNCF Jaeger Distributed Tracing Platform - IAMGuarded variant",
			Expect: true,
		},
		{
			Name:   "volsync",
			Input:  "Asynchronous data replication for Kubernetes volumes",
			Expect: true,
		},
		{
			Name:   "falco-exporter",
			Input:  "Prometheus Metrics Exporter for Falco output events",
			Expect: true,
		},
		{
			Name:   "chrony_exporter-fips",
			Input:  "FIPS-enabled Wolfi-based image for Chrony Exporter - a Prometheus exporter for Chrony NTP metrics.",
			Expect: true,
		},
		{
			Name:   "electric",
			Input:  "Real-time sync for Postgres.",
			Expect: true,
		},
		{
			Name:   "whereabouts",
			Input:  "Whereabouts is a simple IPAM (IP Address Management) solution for Kubernetes. To get more information about Whereabouts, please visit the [official project repository](https://github.com/k8snetworkplumbingwg/whereabouts).",
			Expect: true,
		},
		{
			Name:   "dependency-track",
			Input:  "[Dependency Track](https://github.com/DependencyTrack/dependency-track) Dependency-Track is an intelligent Component Analysis platform that allows organizations to identify and reduce risk in the software supply chain.",
			Expect: true,
		},
		{
			Name:   "chart-testing",
			Input:  "Tool for testing Helm charts, used for linting and testing pull requests.",
			Expect: true,
		},
		{
			Name:   "keycloak-iamguarded-fips",
			Input:  "Minimalist Wolfi-based [Keycloak](https://www.keycloak.org/) IAMGuarded image for identity and access management.",
			Expect: true,
		},
		{
			Name:   "infinispan",
			Input:  "Infinispan is a distributed cache [1] and key-value NoSQL in-memory database",
			Expect: true,
		},
		{
			Name:   "cassandra-medusa",
			Input:  "[cassandra-medusa](https://github.com/thelastpickle/cassandra-medusa), is a Apache Cassandra Backup and Restore Tool.",
			Expect: true,
		},
		{
			Name:   "prometheus-mysqld-exporter",
			Input:  "Minimal Prometheus mysqld exporter Image",
			Expect: true,
		},
		{
			Name:   "envoy-gateway",
			Input:  "Manages Envoy Proxy as a Standalone or Kubernetes-based Application Gateway.",
			Expect: true,
		},
		{
			Name:   "apache-hop",
			Input:  "Data orchestration and engineering platform for managing ETL/ELT workflows.",
			Expect: true,
		},
		{
			Name:   "ntia-conformance-checker",
			Input:  "Check SPDX SBOM for NTIA minimum elements",
			Expect: true,
		},
		{
			Name:   "hubble-export-stdout",
			Input:  "hubble-export-stdout exports Hubble data to stdout.",
			Expect: true,
		},
		{
			Name:   "chainguard-base",
			Input:  "Minimal image useful as a base for building secure images.",
			Expect: true,
		},
		{
			Name:   "ceph-csi-operator",
			Input:  "Operator for Ceph CSI driver management in Kubernetes",
			Expect: true,
		},
		{
			Name:   "pgpool2_exporter-fips",
			Input:  "A FIPS-compliant Prometheus exporter image for Pgpool-II metrics.",
			Expect: true,
		},
		{
			Name:   "fluentd-kubernetes-daemonset-fips",
			Input:  "Wolfi-based images that provide Fluentd DaemonSets for Kubernetes.",
			Expect: true,
		},
		{
			Name:   "envoy-ratelimit",
			Input:  " Go/gRPC service designed to enable generic rate limit scenarios from different types of applications.",
			Expect: true,
		},
		{
			Name:   "sealed-secrets-kubeseal",
			Input:  "A Kubernetes tool that uses one-way encryption of Secrets, enabling safe GitOps-friendly secret management.",
			Expect: true,
		},
		{
			Name:   "rke2-runtime",
			Input:  "Minimal image of RKE2's container runtime",
			Expect: true,
		},
		{
			Name:   "request-7258",
			Input:  "null",
			Expect: true,
		},
		{
			Name:   "dex",
			Input:  "[dex](https://dexidp.io) is a federated OpenID Connect provider.",
			Expect: true,
		},
		{
			Name:   "opentelemetry-collector-contrib-fips",
			Input:  "Minimal image with [opentelemetry-collector-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib).",
			Expect: true,
		},
		{
			Name:   "grpc-health-probe-fips",
			Input:  "A tool to perform health-checks for gRPC applications in Kubernetes and elsewhere",
			Expect: true,
		},
		{
			Name:   "sql_exporter-fips",
			Input:  "Database-agnostic SQL Exporter for Prometheus",
			Expect: true,
		},
		{
			Name:   "cassandra-iamguarded",
			Input:  "Apache Cassandra distributed NoSQL database, compatible with iamguarded charts",
			Expect: true,
		},
		{
			Name:   "cluster-api-gcp-controller",
			Input:  "Kubernetes Cluster API provider for Google Cloud Platform infrastructure management.",
			Expect: true,
		},
		{
			Name:   "rancher-security-scan",
			Input:  "Evaluates Kubernetes cluster security posture against established best practices using kube-bench framework.",
			Expect: true,
		},
		{
			Name:   "contour-iamguarded",
			Input:  "Contour is an ingress controller for Kubernetes that works by deploying the Envoy proxy as a reverse proxy and load balancer. Contour supports dynamic configuration updates out of the box while maintaining a lightweight profile.",
			Expect: true,
		},
		{
			Name:   "jitsucom-jitsu",
			Input:  "Jitsu is an open-source Segment alternative. Fully-scriptable data ingestion engine for modern data teams. Set-up a real-time data pipeline in minutes, not days",
			Expect: true,
		},
		{
			Name:   "fluentd-kubernetes-daemonset",
			Input:  "Wolfi-based images that provide Fluentd DaemonSets for Kubernetes.",
			Expect: true,
		},
		{
			Name:   "apko",
			Input:  "Container image for running [apko](https://github.com/chainguard-dev/apko) container builds.",
			Expect: true,
		},
		{
			Name:   "prometheus",
			Input:  "Chainguard image for Prometheus, a systems and service monitoring system.",
			Expect: true,
		},
		{
			Name:   "onepassword-operator",
			Input:  "The 1Password Connect Kubernetes Operator provides the ability to integrate Kubernetes Secrets with 1Password",
			Expect: true,
		},
		{
			Name:   "trust-manager",
			Input:  "Minimalist Wolfi-based trust-manager operator for distributing trust bundles across a Kubernetes cluster.",
			Expect: true,
		},
		{
			Name:   "unbound-mailcow-fips",
			Input:  "Unbound is a validating, recursive, and caching DNS resolver.",
			Expect: true,
		},
		{
			Name:   "nats-box-fips",
			Input:  "A lightweight container with NATS utilities(FIPS Compliant).",
			Expect: true,
		},
		{
			Name:   "ingress-nginx-controller",
			Input:  " Ingress-NGINX Controller for Kubernetes",
			Expect: true,
		},
		{
			Name:   "wiremock",
			Input:  "Wiremock is a tool for mocking HTTP services.",
			Expect: true,
		},
		{
			Name:   "wave",
			Input:  "Wave watches Deployments within a Kubernetes cluster and ensures that each Deployment's Pods always have up to date configuration.",
			Expect: true,
		},
		{
			Name:   "localstack",
			Input:  "A fully functional local AWS cloud stack. Develop and test your cloud & Serverless apps offline",
			Expect: true,
		},
		{
			Name:   "crossplane-aws-provider",
			Input:  "Crossplane provider-aws is the infrastructure provider for Amazon Web Services (AWS).",
			Expect: true,
		},
		{
			Name:   "victoriametrics-operator-config-reloader",
			Input:  "Config reloader component for VictoriaMetrics Operator that watches and reloads configurations for VMAgent, VMAuth, and VMAlert.",
			Expect: true,
		},
		{
			Name:   "redis-operator",
			Input:  "Redis Operator is a Kubernetes operator image that automates Redis cluster deployment, scaling, and management in Kubernetes environments.",
			Expect: true,
		},
		{
			Name:   "contour-iamguarded-fips",
			Input:  "Contour is an ingress controller for Kubernetes that works by deploying the Envoy proxy as a reverse proxy and load balancer. Contour supports dynamic configuration updates out of the box while maintaining a lightweight profile.",
			Expect: true,
		},
		{
			Name:   "cert-manager-openshift-routes-fips",
			Input:  "A FIPS-compliant image for cert-manager OpenShift Route support",
			Expect: true,
		},
		{
			Name:   "prefect",
			Input:  "A minimal, wolfi-based image for Prefect. Prefect is a modern workflow orchestration framework for building, observing, and reacting to data pipelines.",
			Expect: true,
		},
		{
			Name:   "graalvm-native",
			Input:  "Base image with just enough files to run native [GraalVM](https://www.graalvm.org/) native-image binaries.",
			Expect: true,
		},
		{
			Name:   "dask-gateway",
			Input:  "A multi-tenant server for securely deploying and managing [Dask clusters](https://gateway.dask.org/).",
			Expect: true,
		},
		{
			Name:   "crossplane-aws",
			Input:  "Crossplane provider for managing Amazon Web Services (AWS) config services in Kubernetes.",
			Expect: true,
		},
		{
			Name:   "prometheus-elasticsearch-exporter",
			Input:  "Minimalist Wolfi-based Prometheus Elasticsearch Exporter image for exporting various metrics about Elasticsearch.",
			Expect: true,
		},
		{
			Name:   "sealed-secrets-controller-iamguarded-fips",
			Input:  "A Kubernetes controller and tool for one-way encrypted Secrets, enabling safe GitOps-friendly secret management.",
			Expect: true,
		},
		{
			Name:   "mysql-iamguarded",
			Input:  "MySQL is a widely used open-source relational database management system.",
			Expect: true,
		},
		{
			Name:   "dcgm-exporter-fips",
			Input:  "NVIDIA GPU metrics exporter for Prometheus leveraging DCGM",
			Expect: true,
		},
		{
			Name:   "velero-plugin-for-aws",
			Input:  "Plugins to support Velero on AWS",
			Expect: true,
		},
		{
			Name:   "coredns-fips",
			Input:  "A minimal FIPS image of CoreDNS for secure and flexible DNS-based service discovery in kubernetes environments.",
			Expect: true,
		},
		{
			Name:   "vault-csi-provider",
			Input:  "HashiCorp Vault Provider for Secret Store CSI Driver",
			Expect: true,
		},
		{
			Name:   "kube-fluentd-operator",
			Input:  "This image is used for the [Kubernetes Fluentd Operator](https://github.com/vmware/kube-fluentd-operator)",
			Expect: true,
		},
		{
			Name:   "prometheus-blackbox-exporter-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "pulumi",
			Input:  "Minimal Pulumi Image",
			Expect: true,
		},
		{
			Name:   "pytorch",
			Input:  "A minimal, [wolfi](https://github.com/wolfi-dev)-based image for pytorch, a Python package that provides two high-level features: Tensor computation with strong GPU acceleration and Deep neural networks built on a tape-based autograd system.",
			Expect: true,
		},
		{
			Name:   "kiali-operator",
			Input:  "Kiali Operator that manages the lifecycle of Kiali in Kubernetes environments with Istio service mesh integration",
			Expect: true,
		},
		{
			Name:   "opensearch-dashboards",
			Input:  "Minimal image with OpenSearch Dashboards",
			Expect: true,
		},
		{
			Name:   "socat",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "http-echo",
			Input:  "http-echo is a lightweight Go-based web server that responds to all HTTP requests with a predefined message, making it useful for testing and debugging HTTP interactions.",
			Expect: true,
		},
		{
			Name:   "argocd-image-updater-fips",
			Input:  "Automatic container image update for Argo CD",
			Expect: true,
		},
		{
			Name:   "pg-timetable",
			Input:  "An advanced standalone job scheduler for PostgreSQL, offering many advantages over traditional schedulers such as cron and others.",
			Expect: true,
		},
		{
			Name:   "openldap-fips",
			Input:  "OpenLDAP is a free, open-source implementation of the Lightweight Directory Access Protocol (LDAP) developed by the OpenLDAP Project.",
			Expect: true,
		},
		{
			Name:   "rancher-fleet-fips",
			Input:  "Deploy workloads from Git to large fleets of Kubernetes clusters",
			Expect: true,
		},
		{
			Name:   "google-cloud-sdk-iamguarded-fips",
			Input:  "Minimal FIPS-compliant IAMGuarded image with the [Google Cloud SDK](https://cloud.google.com/sdk/).",
			Expect: true,
		},
		{
			Name:   "cluster-api-azure-controller",
			Input:  "Kubernetes Cluster API provider for Microsoft Azure infrastructure management.",
			Expect: true,
		},
		{
			Name:   "opencost",
			Input:  "OpenCost give teams visibility into current and historical Kubernetes and cloud spend and resource allocation.",
			Expect: true,
		},
		{
			Name:   "volcano-fips",
			Input:  "A Kubernetes-native batch scheduling system, extending and enhancing the capabilities of the standard kube-scheduler.",
			Expect: true,
		},
		{
			Name:   "argo-fips",
			Input:  "Argo is a collection of tools for Kubernetes that help users to run workflows and manage clusters.",
			Expect: true,
		},
		{
			Name:   "caddy",
			Input:  "Open source web server with automatic HTTPS written in Go",
			Expect: true,
		},
		{
			Name:   "spicedb-fips",
			Input:  "This is the FIPS-compliant variant of Spice db an open-source authorization database inspired by Google's Zanzibar, providing scalable and fine-grained access control for applications.",
			Expect: true,
		},
		{
			Name:   "jenkins",
			Input:  "A minimal, Wolfi-based container image for Jenkins - an open-source CI/CD server that enables developers to build, test, and deploy their software.",
			Expect: true,
		},
		{
			Name:   "keycloak-operator",
			Input:  "A Kubernetes Operator based on the Operator SDK for installing and managing Keycloak.",
			Expect: true,
		},
		{
			Name:   "node",
			Input:  "Minimal container image for running NodeJS apps",
			Expect: true,
		},
		{
			Name:   "kubernetes-ingress-defaultbackend-fips",
			Input:  "Minimal image that acts as a drop-in replacement for the `registry.k8s.io/defaultbackend` image. Used in some ingresses like https://github.com/kubernetes/ingress-gce and https://github.com/kubernetes/ingress-nginx",
			Expect: true,
		},
		{
			Name:   "grafana-pyroscope",
			Input:  "[Grafana Pyroscope](https://grafana.com/oss/pyroscope/) is a continuous profiling platform that allows you to debug performance issues down to a single line of code.",
			Expect: true,
		},
		{
			Name:   "opentelemetry-operator-fips",
			Input:  "Kubernetes Operator for OpenTelemetry Collector",
			Expect: true,
		},
		{
			Name:   "headlamp-plugin-flux-fips",
			Input:  "Headlamp plugin to visualize and manage Flux GitOps resources in Kubernetes clusters.",
			Expect: true,
		},
		{
			Name:   "wave-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "docker-compose",
			Input:  "minimal docker-compose image with docker-compose binary",
			Expect: true,
		},
		{
			Name:   "apisix-ingress-controller-iamguarded",
			Input:  "Apache APISIX Ingress Controller for Kubernetes ingress management.",
			Expect: true,
		},
		{
			Name:   "sqlpad",
			Input:  "A minimal Wolfi-based image for sqlpad, which is a web application for generating and running SQL queries and visualizing the results. For more information, please refer to the [applications documentation](https://github.com/sqlpad/sqlpad) on github.",
			Expect: true,
		},
		{
			Name:   "request-5792",
			Input:  "The JuiceFS Container Storage Interface (CSI) Driver Image that implements the CSI specification for container orchestrators to manage the lifecycle of JuiceFS filesystems.",
			Expect: true,
		},
		{
			Name:   "sonarqube",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "gatus",
			Input:  "Gatus is a dev-oriented health dashboard that gives you the ability to monitor your services using HTTP, ICMP, TCP and DNS queries",
			Expect: true,
		},
		{
			Name:   "newrelic-k8s-events-forwarder",
			Input:  "Minimal [newrelic-k8s-events-forwarder](https://github.com/newrelic/nri-kubernetes) container image.",
			Expect: true,
		},
		{
			Name:   "victoriametrics-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "cephcsi-fips",
			Input:  "CephCSI is the Container Storage Interface (CSI) driver for Ceph, providing support for RBD and CephFS.",
			Expect: true,
		},
		{
			Name:   "unbound",
			Input:  "Unbound is a validating, recursive, and caching DNS resolver.",
			Expect: true,
		},
		{
			Name:   "x509-certificate-exporter",
			Input:  "A Prometheus exporter to monitor x509 certificates expiration in Kubernetes clusters or standalone",
			Expect: true,
		},
		{
			Name:   "cilium-envoy",
			Input:  "Cilium-envoy is a specialized Envoy proxy used by Cilium for Layer 7 policy enforcement and service mesh functionality.",
			Expect: true,
		},
		{
			Name:   "cluster-api-helm-controller-fips",
			Input:  "CAAPH uses Helm charts to manage the installation and lifecycle of Cluster API add-ons.",
			Expect: true,
		},
		{
			Name:   "pgbouncer",
			Input:  "This image contains the CLI for the [pgbouncer](https://www.pgbouncer.org/) connection pooler for PostgreSQL. This image contains the `pgbouncer` binary and can be used directly.",
			Expect: true,
		},
		{
			Name:   "strimzi-kafka-operator",
			Input:  "Strimzi provides a way to run an Apache Kafka cluster on Kubernetes in various deployment configurations.",
			Expect: true,
		},
		{
			Name:   "jupyterhub-k8s-image-awaiter",
			Input:  "JupyterHub Kubernetes Image Awaiter - ensures images are pre-pulled before deployment",
			Expect: true,
		},
		{
			Name:   "memcached-exporter",
			Input:  "A memcached exporter for Prometheus.",
			Expect: true,
		},
		{
			Name:   "spark-operator-fips",
			Input:  "A minimal, FIPS 140-3 compliant image for Spark Operator. Facilitates the deployment and management of Apache Spark applications in Kubernetes environments.",
			Expect: true,
		},
		{
			Name:   "management-api-for-apache-cassandra",
			Input:  "RESTful / Secure Management Sidecar for Apache Cassandra",
			Expect: true,
		},
		{
			Name:   "zookeeper-iamguarded-fips",
			Input:  "[Apache ZooKeeper](https://zookeeper.apache.org/) is an effort to develop and maintain an open-source server which enables highly reliable distributed coordination.",
			Expect: true,
		},
		{
			Name:   "seaweedfs",
			Input:  "SeaweedFS is a fast distributed storage system for blobs, objects, files, and data lake, providing S3-compatible API and filesystem interface",
			Expect: true,
		},
		{
			Name:   "kafka-proxy-fips",
			Input:  "Proxy connections to Kafka cluster. Connect through SOCKS Proxy, HTTP Proxy or to cluster running in Kubernetes.",
			Expect: true,
		},
		{
			Name:   "prometheus-iamguarded-fips",
			Input:  "Prometheus is a monitoring system and time series database",
			Expect: true,
		},
		{
			Name:   "kubernetes",
			Input:  "Production-Grade Container Scheduling and Management",
			Expect: true,
		},
		{
			Name:   "coredns",
			Input:  "A minimal image of CoreDNS for secure and flexible DNS-based service discovery in kubernetes environments.",
			Expect: true,
		},
		{
			Name:   "ollama",
			Input:  "Get up and running with Llama 3.3, DeepSeek-R1, Phi-4, Gemma 3, Mistral Small 3.1 and other large language models.",
			Expect: true,
		},
		{
			Name:   "secrets-store-csi-driver-crds-fips",
			Input:  "An image to be used as a init container to create the necessary CRDs while deploying the Kubernetes Secrets Store CSI Driver.",
			Expect: true,
		},
		{
			Name:   "thanos-iamguarded-fips",
			Input:  "Highly available Prometheus setup with long term storage",
			Expect: true,
		},
		{
			Name:   "datadog-operator-fips",
			Input:  "Kubernetes Operator for Datadog Resources",
			Expect: true,
		},
		{
			Name:   "TODO find out name",
			Input:  "Advanced object-relational database management system",
			Expect: true,
		},
		{
			Name:   "liquibase",
			Input:  "Liquibase is a database schema change management solution that enables you to revise and release database changes faster and safer from development to production.",
			Expect: true,
		},
		{
			Name:   "go-openssl",
			Input:  "Golang toolchain with golang-fips/go patchset.",
			Expect: true,
		},
		{
			Name:   "aws-network-policy-agent",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "nats-fips",
			Input:  "NATS is a flexible messaging system providing pub/sub, streaming, storage etc.",
			Expect: true,
		},
		{
			Name:   "newrelic-infrastructure-bundle",
			Input:  "Minimal [newrelic-infrastructure-bundle](https://github.com/newrelic/infrastructure-bundle) container image.",
			Expect: true,
		},
		{
			Name:   "wso2is",
			Input:  "WSO2 Identity Server is a powerful, modern identity and access management solution for your on-premises or cloud environment",
			Expect: true,
		},
		{
			Name:   "secrets-store-csi-driver-provider-aws",
			Input:  "The AWS provider for the Secrets Store CSI Driver allows you to fetch secrets from AWS Secrets Manager and AWS Systems Manager Parameter Store, and mount them into Kubernetes pods.",
			Expect: true,
		},
		{
			Name:   "proxysql",
			Input:  "Minimal image with [proxysql](https://github.com/sysown/proxysql).",
			Expect: true,
		},
		{
			Name:   "prometheus-operator-iamguarded-fips",
			Input:  "Prometheus Operator creates/configures/manages Prometheus clusters atop Kubernetes",
			Expect: true,
		},
		{
			Name:   "external-secrets",
			Input:  "Fetches secrets from external systems and exposes them as Kubernetes Secrets.",
			Expect: true,
		},
		{
			Name:   "jaeger-operator",
			Input:  "Minimal jaeger-operator container image.",
			Expect: true,
		},
		{
			Name:   "kubernetes-autoscaler-addon-resizer",
			Input:  "Addon-resizer is a container that vertically scales a Deployment based on the number of nodes in your cluster.",
			Expect: true,
		},
		{
			Name:   "cilium-envoy-fips",
			Input:  "Cilium-envoy-fips is a FIPS 140-3 compliant, specialized Envoy proxy used by Cilium for Layer 7 policy enforcement and service mesh functionality.",
			Expect: true,
		},
		{
			Name:   "nats-server-config-reloader",
			Input:  "Monitors NATS configuration files and triggers reloads without restarting the server.",
			Expect: true,
		},
		{
			Name:   "nginx-s3-gateway-fips",
			Input:  "NGINX S3 Gateway",
			Expect: true,
		},
		{
			Name:   "memcached-exporter-fips",
			Input:  "A memcached exporter for Prometheus.",
			Expect: true,
		},
		{
			Name:   "cert-manager-openshift-routes",
			Input:  "OpenShift Route support for cert-manager",
			Expect: true,
		},
		{
			Name:   "cluster-api-aws-controller",
			Input:  "Kubernetes Cluster API Provider AWS provides consistent deployment and day 2 operations of \"self-managed\" and EKS Kubernetes clusters on AWS",
			Expect: true,
		},
		{
			Name:   "cert-manager-webhook-pdns-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "keycloak",
			Input:  "Minimalist Wolfi-based [Keycloak](https://www.keycloak.org/) image for identity and access management.",
			Expect: true,
		},
		{
			Name:   "helm",
			Input:  "Minimal image with [helm](https://helm.sh) binary.",
			Expect: true,
		},
		{
			Name:   "newrelic-fluent-bit-output-fips",
			Input:  "Minimal [newrelic-fluent-bit-output](https://github.com/newrelic/newrelic-fluent-bit-output) container image.",
			Expect: true,
		},
		{
			Name:   "aws-eks-pod-identity-agent",
			Input:  "EKS Pod Identity is a feature of Amazon EKS that simplifies the process for cluster administrators to configure Kubernetes applications with AWS IAM permissions",
			Expect: true,
		},
		{
			Name:   "kubeflow",
			Input:  "Kubeflow is a Machine Learning Toolkit for Kubernetes",
			Expect: true,
		},
		{
			Name:   "kubernetes-csi-external-resizer",
			Input:  "Minimal image with [kubernetes-csi/external-resizer](https://github.com/kubernetes-csi/external-resizer).",
			Expect: true,
		},
		{
			Name:   "os-shell-iamguarded",
			Input:  "OS Shell + Utility is a general-purpose minimal image, used by Iamguarded Helm Charts.",
			Expect: true,
		},
		{
			Name:   "jitsucom-bulker",
			Input:  "Service for bulk-loading data to databases with automatic schema management (Redshift, Snowflake, BigQuery, ClickHouse, Postgres, MySQL)",
			Expect: true,
		},
		{
			Name:   "chisel-fips",
			Input:  "A fast TCP/UDP tunnel over HTTP",
			Expect: true,
		},
		{
			Name:   "tritonserver-no-backend-fips",
			Input:  "The Triton Inference Server provides an optimized cloud and edge inferencing solution.",
			Expect: true,
		},
		{
			Name:   "copybara-fips",
			Input:  "FIPS-compliant Copybara image for transforming and moving code between repositories, enabling code synchronization workflows between different version control systems with FIPS 140-2/140-3 validated cryptography.",
			Expect: true,
		},
		{
			Name:   "cis-operator",
			Input:  "Enables running CIS benchmark security scans on a Kubernetes cluster and generates compliance reports that can be downloaded",
			Expect: true,
		},
		{
			Name:   "request-5190",
			Input:  "A minimal, [wolfi](https://github.com/wolfi-dev)-based image for pytorch, a Python package that provides two high-level features: Tensor computation with strong GPU acceleration and Deep neural networks built on a tape-based autograd system.",
			Expect: true,
		},
		{
			Name:   "request-4089",
			Input:  "A minimal, [wolfi](https://github.com/wolfi-dev)-based image for vllm, a high-throughput and memory-efficient inference and serving engine for LLMs",
			Expect: true,
		},
		{
			Name:   "openldap",
			Input:  "OpenLDAP is a free, open-source implementation of the Lightweight Directory Access Protocol (LDAP) developed by the OpenLDAP Project.",
			Expect: true,
		},
		{
			Name:   "argo",
			Input:  "Argo is a collection of tools for Kubernetes that help users to run workflows and manage clusters.",
			Expect: true,
		},
		{
			Name:   "cilium",
			Input:  "[Cilium](https://cilium.io/) is an open source, cloud native solution for providing, securing, and observing network connectivity between workloads using eBPF.",
			Expect: true,
		},
		{
			Name:   "aws-sigv4-proxy",
			Input:  "This project signs and proxies HTTP requests with Sigv4",
			Expect: true,
		},
		{
			Name:   "clickhouse-operator-fips",
			Input:  "Kubernetes Operator for ClickHouse. Creates, configures and manages ClickHouse clusters running on Kubernetes.",
			Expect: true,
		},
		{
			Name:   "proxysql-fips",
			Input:  "Minimal image with [proxysql](https://github.com/sysown/proxysql).",
			Expect: true,
		},
		{
			Name:   "kyverno-policy-reporter-plugin-kyverno-fips",
			Input:  "This Plugin for Policy Reporter brings additional Kyverno specific information to the Policy Reporter UI",
			Expect: true,
		},
		{
			Name:   "kubectl",
			Input:  "Minimal image with kubectl binary.",
			Expect: true,
		},
		{
			Name:   "cloud-provider-gcp-cloud-controller-manager",
			Input:  "Kubernetes cloud controller manager for Google Cloud Platform (GCP), managing cloud-specific resources and integrations.",
			Expect: true,
		},
		{
			Name:   "tritonserver-pytorch-backend-fips",
			Input:  "The Triton backend for the PyTorch TorchScript models with FIPS support.",
			Expect: true,
		},
		{
			Name:   "pushprox",
			Input:  "PushProx is a client and proxy that allows transversing of NAT and other similar network topologies by Prometheus",
			Expect: true,
		},
		{
			Name:   "gatekeeper-crds-fips",
			Input:  "Minimal FIPS-compliant image for installing Gatekeeper Custom Resource Definitions (CRDs) in Kubernetes clusters.",
			Expect: true,
		},
		{
			Name:   "k8s-agents-operator-fips",
			Input:  "k8s-agents-operator auto-instruments containerized workloads in Kubernetes with New Relic agents.",
			Expect: true,
		},
		{
			Name:   "it-tools",
			Input:  "`it-tools` is a collection of useful tools for developer and people working in IT.",
			Expect: true,
		},
		{
			Name:   "zitadel",
			Input:  "ZITADEL is an open-source identity and access management (IAM) system that simplifies user authentication and authorization for applications.",
			Expect: true,
		},
		{
			Name:   "flux-source-watcher-fips",
			Input:  "FIPS-validated Flux Source Watcher extends Flux CD with ArtifactGenerator CRD for source composition",
			Expect: true,
		},
		{
			Name:   "opentelemetry-java-instrumentation",
			Input:  "OpenTelemetry auto-instrumentation and instrumentation libraries for Java",
			Expect: true,
		},
		{
			Name:   "clamav",
			Input:  "ClamAV® is an open source antivirus engine for detecting trojans, viruses, malware & other malicious threats.",
			Expect: true,
		},
		{
			Name:   "kubernetes-dashboard-api",
			Input:  "Stateless Go module, which could be referred to as a Kubernetes API extension",
			Expect: true,
		},
		{
			Name:   "weaviate",
			Input:  "Minimal container image for running the weaviate vector database.",
			Expect: true,
		},
		{
			Name:   "garage",
			Input:  "Garage is an S3-compatible distributed object storage service.",
			Expect: true,
		},
		{
			Name:   "thanos-receive-controller",
			Input:  "Kubernetes controller to automatically configure Thanos receive hashrings",
			Expect: true,
		},
		{
			Name:   "request-4334",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "victoriametrics-operator-config-reloader-fips",
			Input:  "Config reloader component for VictoriaMetrics Operator that watches and reloads configurations for VMAgent, VMAuth, and VMAlert.",
			Expect: true,
		},
		{
			Name:   "openssh-server",
			Input:  "OpenSSH Server is a secure shell (SSH) server implementation that provides encrypted communication between clients and servers.",
			Expect: true,
		},
		{
			Name:   "redis-iamguarded",
			Input:  "[Redis](https://github.com/redis/redis) Redis is an in-memory data structure store, used as a database, cache, and message broker.",
			Expect: true,
		},
		{
			Name:   "kube-arangodb-fips",
			Input:  "FIPS-compliant ArangoDB Kubernetes Operator for managing ArangoDB database deployments",
			Expect: true,
		},
		{
			Name:   "kyverno",
			Input:  "[Kyverno](https://kyverno.io/) is a policy engine that allows you to write policies as Kubernetes resources and manage them with familiar tools",
			Expect: true,
		},
		{
			Name:   "kubeflow-pipelines",
			Input:  "Minimalist Kubeflow Pipelines Images",
			Expect: true,
		},
		{
			Name:   "dex-iamguarded",
			Input:  "dex is a federated OpenID Connect provider.",
			Expect: true,
		},
		{
			Name:   "kafka-exporter",
			Input:  "Kafka exporter for Prometheus",
			Expect: true,
		},
		{
			Name:   "postgres-operator",
			Input:  "Creates and manages PostgreSQL clusters running in Kubernetes.",
			Expect: true,
		},
		{
			Name:   "pgpool2-fips",
			Input:  "Middleware that works between PostgreSQL servers and a PostgreSQL database client.",
			Expect: true,
		},
		{
			Name:   "kor",
			Input:  "A Golang Tool to discover unused Kubernetes Resources",
			Expect: true,
		},
		{
			Name:   "cosign-fips",
			Input:  "Minimalist Wolfi-based Cosign image for signing and verifying images using Sigstore.",
			Expect: true,
		},
		{
			Name:   "os-shell-iamguarded-fips",
			Input:  "OS Shell + Utility is a general-purpose minimal image, used by Iamguarded Helm Charts.",
			Expect: true,
		},
		{
			Name:   "grype-fips",
			Input:  "A vulnerability scanner for container images and filesystems. This image is built and tested for FIPS compliance to meet strict security standards in regulated environments.",
			Expect: true,
		},
		{
			Name:   "apisix-ingress-controller",
			Input:  "Minimal wolfi image of Apache APISIX Ingress Controller use to run APISIX Gateway as a Kubernetes Ingress to handle inbound traffic for a Kubernetes cluster.",
			Expect: true,
		},
		{
			Name:   "request-5011",
			Input:  "A minimal, Wolfi-based image for Spark Operator. Facilitates the deployment and management of Apache Spark applications in Kubernetes environments.",
			Expect: true,
		},
		{
			Name:   "prometheus-config-reloader",
			Input:  "Minimalist image for Prometheus Config Reloader. It helps with config of Prometheus Operator which creates/configures/manages Prometheus clusters atop Kubernetes",
			Expect: true,
		},
		{
			Name:   "memcached-iamguarded",
			Input:  "[Memcached](https://memcached.org/) is an in-memory key-value store for small chunks of arbitrary data (strings, objects) from results of database calls, API calls, or page rendering.",
			Expect: true,
		},
		{
			Name:   "ztunnel-fips",
			Input:  "The ztunnel component of ambient mesh",
			Expect: true,
		},
		{
			Name:   "prometheus-statsd-exporter",
			Input:  "Minimalist Wolfi-based Prometheus StatsD Exporter image for exporting metrics to StatsD.",
			Expect: true,
		},
		{
			Name:   "secrets-store-csi-driver-crds",
			Input:  "An image to be used as a init container to create the necessary CRDs while deploying the Kubernetes Secrets Store CSI Driver.",
			Expect: true,
		},
		{
			Name:   "volcano",
			Input:  "A Kubernetes-native batch scheduling system, extending and enhancing the capabilities of the standard kube-scheduler.",
			Expect: true,
		},
		{
			Name:   "vector",
			Input:  "Minimal image with [Vector](https://vector.dev/), an end-to-end data observability pipeline",
			Expect: true,
		},
		{
			Name:   "imagemagick",
			Input:  "ImageMagick® is a free and open-source software suite, used for editing and manipulating digital images.",
			Expect: true,
		},
		{
			Name:   "kube-logging-operator-fluentd",
			Input:  "Kubernetes Logging Operator Fluentd Image",
			Expect: true,
		},
		{
			Name:   "cedar",
			Input:  "This image contains the CLI for the [Cedar Policy](https://www.cedarpolicy.com/en) Language. The binary can be used to run, test, format, or evaluate Cedar policies.",
			Expect: true,
		},
		{
			Name:   "ffmpeg-fips",
			Input:  "Minimal image that contains ffmpeg",
			Expect: true,
		},
		{
			Name:   "aws-for-fluent-bit-fips",
			Input:  "",
			Expect: true,
		},
		{
			Name:   "pdns-auth-fips",
			Input:  "PowerDNS Authoritative Server - High-performance DNS server with flexible backend support.",
			Expect: true,
		},
		{
			Name:   "amazon-cloudwatch-agent-fips",
			Input:  "CloudWatch Agent enables you to collect and export host-level metrics and logs on instances running Linux or Windows server. ",
			Expect: true,
		},
		{
			Name:   "kube-arangodb",
			Input:  "ArangoDB Kubernetes Operator for managing ArangoDB database deployments",
			Expect: true,
		},
		{
			Name:   "harbor-exporter",
			Input:  "A Wolf-based image for Harbor Exporter - application for monitoring harbor deployments.",
			Expect: true,
		},
		{
			Name:   "crossplane-function-go-templating",
			Input:  "This composition function allows you to compose Crossplane resources using Go templates.",
			Expect: true,
		},
		{
			Name:   "crossplane-function-environment-configs",
			Input:  "Crossplane function that manages environment-specific configurations for resources in compositions",
			Expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := ValidateDescription(tt.Input)
			if (got == nil) != tt.Expect {
				t.Errorf("Expected (`%s`) to return (err == nil) == %v, but got %v", tt.Input, tt.Expect, got)
			}
		})
	}
}
