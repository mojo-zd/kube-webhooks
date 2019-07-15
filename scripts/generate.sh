dist_dir=../deployments/
template_dir=./templates/
admission_output_name=admissionregistration.yaml
secret_output_name=secret.yaml

# clear file
rm -rf ${dist_dir}${admission_output_name}
rm -rf ${dist_dir}${secret_output_name}

configuration=`cat ${template_dir}admissionregistration-template.yaml`
secret=`cat ${template_dir}secret-template.yaml`

export caBundle=`cat ca.crt | base64`
export serverCrt=`cat server.crt | base64`
export serverKey=`cat server.key | base64`

echo "${configuration}" | mo  >> ${dist_dir}${admission_output_name}
echo "${secret}" | mo  >> ${dist_dir}${secret_output_name}
