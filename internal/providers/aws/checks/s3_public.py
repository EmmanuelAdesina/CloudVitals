#!/usr/bin/env python3
import argparse
import json
import boto3
from botocore.exceptions import ClientError


def check(profile, region):
    try:
        session = boto3.Session(profile_name=profile, region_name=region)
        s3 = session.client("s3")

        findings = []
        buckets = s3.list_buckets().get("Buckets", [])

        for bucket in buckets:
            name = bucket["Name"]
            try:
                resp = s3.get_public_access_block(Bucket=name)
                config = resp.get("PublicAccessBlockConfiguration", {})

                if not all([
                    config.get("BlockPublicAcls", False),
                    config.get("IgnorePublicAcls", False),
                    config.get("BlockPublicPolicy", False),
                    config.get("RestrictPublicBuckets", False),
                ]):
                    findings.append({
                        "resource": f"arn:aws:s3:::{name}",
                        "region": region,
                        "detail": "S3 bucket public access block not fully enabled",
                        "remediation": (
                            f"aws s3api put-public-access-block --bucket {name} "
                            f"--public-access-block-configuration "
                            f"BlockPublicAcls=true,IgnorePublicAcls=true,"
                            f"BlockPublicPolicy=true,RestrictPublicBuckets=true"
                        )
                    })
            except ClientError:
                continue

        status = "fail" if findings else "pass"

        print(json.dumps({
            "check_id": "s3_public",
            "check_name": "Public S3 Bucket Access",
            "status": status,
            "severity": "critical",
            "findings": findings,
            "executed_at": None,
            "error_msg": None
        }))

    except Exception as e:
        print(json.dumps({
            "check_id": "s3_public",
            "check_name": "Public S3 Bucket Access",
            "status": "error",
            "severity": "critical",
            "findings": [],
            "executed_at": None,
            "error_msg": str(e)
        }))


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", default="default")
    parser.add_argument("--region", default="us-east-1")
    args = parser.parse_args()
    check(args.profile, args.region)
